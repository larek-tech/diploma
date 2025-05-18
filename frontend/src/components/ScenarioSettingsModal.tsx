import {DomainApiService} from '@/api/DomainApiService';
import {Scenario} from '@/api/models/domain';
import {Button} from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import {Switch} from '@/components/ui/switch';
import {Tabs, TabsContent, TabsList, TabsTrigger} from '@/components/ui/tabs';
import {Textarea} from '@/components/ui/textarea';
import {useToast} from '@/components/ui/use-toast';
import {useStores} from '@/hooks/useStores';
import {Loader2, Plus, Settings} from 'lucide-react';
import {observer} from 'mobx-react-lite';
import {useEffect, useState} from 'react';

interface ScenarioSettingsModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
}

const ScenarioSettingsModal = observer(({ open, onOpenChange }: ScenarioSettingsModalProps) => {
    const { rootStore } = useStores();
    const { toast } = useToast();

    const [selectedScenarioId, setSelectedScenarioId] = useState<number | null>(
        rootStore.selectedScenarioId
    );
    const [availableScenarios, setAvailableScenarios] = useState<Scenario[]>([]);
    const [loadingScenarios, setLoadingScenarios] = useState(false);
    const [activeTab, setActiveTab] = useState('select');

    const [isCreating, setIsCreating] = useState(false);
    const [useMultiquery, setUseMultiquery] = useState(false);
    const [useRerank, setUseRerank] = useState(false);
    const [scenarioData, setScenarioData] = useState<Scenario>({
        id: 0,
        domainId: rootStore.selectedDomainId || 0,
        model: {
            modelName: '',
            systemPrompt: '',
            temperature: 0.7,
            topK: 5,
            topP: 0.9,
        },
        multiQuery: {
            nQueries: 3,
            queryModelName: '',
            useMultiquery: false,
        },
        reranker: {
            rerankerMaxLength: 512,
            rerankerModel: '',
            topK: 5,
            useRerank: false,
        },
        vectorSearch: {
            searchByQuery: true,
            threshold: 0.5,
            topN: 10,
        },
    });

    useEffect(() => {
        if (open) {
            loadScenarios();
        }
    }, [open, rootStore.selectedDomainId]);

    // When switching to create tab, initialize with selected scenario as template if available
    useEffect(() => {
        if (activeTab === 'create' && selectedScenarioId) {
            const selectedScenario = availableScenarios.find((s) => s.id === selectedScenarioId);
            if (selectedScenario) {
                setScenarioData({
                    ...selectedScenario,
                    domainId: rootStore.selectedDomainId || 0,
                });
                setUseMultiquery(selectedScenario.multiQuery.useMultiquery);
                setUseRerank(selectedScenario.reranker.useRerank);
            }
        }
    }, [activeTab, selectedScenarioId, availableScenarios, rootStore.selectedDomainId]);

    const loadScenarios = async () => {
        if (!rootStore.selectedDomainId) return;

        setLoadingScenarios(true);

        try {
            const { scenarios } = await DomainApiService.getScenarios();
            console.log('scenarios', scenarios);

            const allowedScenarioIds = rootStore.selectedDomain?.scenarioIds || [];

            const filteredScenarios = scenarios.filter((scenario) =>
                allowedScenarioIds.includes(scenario.id)
            );

            setAvailableScenarios(filteredScenarios);
        } catch (error) {
            console.error('Ошибка при загрузке домена:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить доступные сценарии',
                variant: 'destructive',
            });
        } finally {
            setLoadingScenarios(false);
        }
    };

    const handleSelectScenario = async (scenarioId: number) => {
        setSelectedScenarioId(scenarioId);
    };

    const handleApplySelectedScenario = async () => {
        if (!selectedScenarioId) {
            toast({
                title: 'Внимание',
                description: 'Выберите сценарий для применения',
                variant: 'destructive',
            });
            return;
        }

        try {
            await rootStore.setSelectedScenario(selectedScenarioId);
            toast({
                title: 'Успех',
                description: 'Сценарий успешно применен',
                variant: 'default',
            });
            onOpenChange(false);
        } catch (error) {
            console.error('Ошибка при применении сценария:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось применить сценарий',
                variant: 'destructive',
            });
        }
    };

    const handleInputChange = (
        section: 'model' | 'multiQuery' | 'reranker' | 'vectorSearch',
        field: string,
        value: string | number | boolean
    ) => {
        setScenarioData((prev) => ({
            ...prev,
            [section]: {
                ...prev[section],
                [field]: value,
            },
        }));
    };

    const handleCreateScenario = async () => {
        if (!rootStore.selectedDomainId) {
            toast({
                title: 'Ошибка',
                description: 'Необходимо выбрать домен',
                variant: 'destructive',
            });
            return;
        }

        setIsCreating(true);
        try {
            // Обновляем значения useMultiquery и useRerank в сценарии
            const newScenario = {
                ...scenarioData,
                domainId: rootStore.selectedDomainId,
                multiQuery: {
                    ...scenarioData.multiQuery,
                    useMultiquery,
                },
                reranker: {
                    ...scenarioData.reranker,
                    useRerank,
                },
            };

            const result = await DomainApiService.createScenario(newScenario);

            // Update the list of scenarios
            await loadScenarios();

            // Select the newly created scenario
            if (result) {
                await rootStore.setSelectedScenario(result.id);
                setSelectedScenarioId(result.id);
            }

            toast({
                title: 'Успех',
                description: 'Сценарий успешно создан и применен',
                variant: 'default',
            });

            // Close the modal
            onOpenChange(false);
        } catch (error) {
            console.error('Ошибка при создании сценария:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось создать сценарий',
                variant: 'destructive',
            });
        } finally {
            setIsCreating(false);
        }
    };

    const renderScenarioDetails = (scenario: Scenario) => {
        return (
            <div className='p-4 mt-4 bg-gray-100 dark:bg-gray-800 rounded-md space-y-4'>
                <div className='space-y-3'>
                    <h4 className='font-medium text-primary'>Настройки модели</h4>
                    <div className='grid grid-cols-2 gap-2 text-sm'>
                        <div>
                            <span className='font-medium'>Модель:</span>
                            <p>{scenario.model.modelName || 'Не указано'}</p>
                        </div>
                        <div>
                            <span className='font-medium'>Температура:</span>
                            <p>{scenario.model.temperature}</p>
                        </div>
                        <div>
                            <span className='font-medium'>Top K:</span>
                            <p>{scenario.model.topK}</p>
                        </div>
                        <div>
                            <span className='font-medium'>Top P:</span>
                            <p>{scenario.model.topP}</p>
                        </div>
                    </div>
                    {scenario.model.systemPrompt && (
                        <div>
                            <span className='font-medium'>Системный промпт:</span>
                            <p className='text-xs mt-1 bg-gray-50 dark:bg-gray-900 p-2 rounded border border-gray-200 dark:border-gray-700 max-h-24 overflow-y-auto'>
                                {scenario.model.systemPrompt}
                            </p>
                        </div>
                    )}
                </div>

                <div className='space-y-3'>
                    <h4 className='font-medium text-primary'>Векторный поиск</h4>
                    <div className='grid grid-cols-2 gap-2 text-sm'>
                        <div>
                            <span className='font-medium'>Поиск по запросу:</span>
                            <p>{scenario.vectorSearch.searchByQuery ? 'Включен' : 'Выключен'}</p>
                        </div>
                        <div>
                            <span className='font-medium'>Порог:</span>
                            <p>{scenario.vectorSearch.threshold}</p>
                        </div>
                        <div>
                            <span className='font-medium'>Top N:</span>
                            <p>{scenario.vectorSearch.topN}</p>
                        </div>
                    </div>
                </div>

                <div className='space-y-3'>
                    <h4 className='font-medium text-primary'>Мультизапрос</h4>
                    <div className='grid grid-cols-2 gap-2 text-sm'>
                        <div>
                            <span className='font-medium'>Статус:</span>
                            <p>{scenario.multiQuery.useMultiquery ? 'Включен' : 'Выключен'}</p>
                        </div>
                        {scenario.multiQuery.useMultiquery && (
                            <>
                                <div>
                                    <span className='font-medium'>Количество запросов:</span>
                                    <p>{scenario.multiQuery.nQueries}</p>
                                </div>
                                <div>
                                    <span className='font-medium'>Модель запроса:</span>
                                    <p>{scenario.multiQuery.queryModelName || 'Не указано'}</p>
                                </div>
                            </>
                        )}
                    </div>
                </div>

                <div className='space-y-3'>
                    <h4 className='font-medium text-primary'>Ранжирование</h4>
                    <div className='grid grid-cols-2 gap-2 text-sm'>
                        <div>
                            <span className='font-medium'>Статус:</span>
                            <p>{scenario.reranker.useRerank ? 'Включено' : 'Выключено'}</p>
                        </div>
                        {scenario.reranker.useRerank && (
                            <>
                                <div>
                                    <span className='font-medium'>Модель ранжирования:</span>
                                    <p>{scenario.reranker.rerankerModel || 'Не указано'}</p>
                                </div>
                                <div>
                                    <span className='font-medium'>Максимальная длина:</span>
                                    <p>{scenario.reranker.rerankerMaxLength}</p>
                                </div>
                                <div>
                                    <span className='font-medium'>Top K:</span>
                                    <p>{scenario.reranker.topK}</p>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </div>
        );
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className='max-w-3xl max-h-[85vh] overflow-y-auto'>
                <DialogHeader>
                    <DialogTitle className='flex items-center gap-2'>
                        <Settings className='h-5 w-5' />
                        Настройки сценария
                    </DialogTitle>
                    <DialogDescription>
                        Выберите существующий сценарий или создайте новый на основе существующего
                    </DialogDescription>
                </DialogHeader>

                <Tabs defaultValue='select' value={activeTab} onValueChange={setActiveTab}>
                    <TabsList className='grid grid-cols-2 mb-4'>
                        <TabsTrigger value='select'>Выбрать сценарий</TabsTrigger>
                        <TabsTrigger value='create'>Создать сценарий</TabsTrigger>
                    </TabsList>

                    <TabsContent value='select' className='space-y-4'>
                        {loadingScenarios ? (
                            <div className='flex justify-center p-8'>
                                <Loader2 className='h-8 w-8 animate-spin text-primary' />
                            </div>
                        ) : availableScenarios.length === 0 ? (
                            <div className='text-center p-8'>
                                <p>Для выбранного домена не найдено сценариев</p>
                                <Button
                                    variant='outline'
                                    onClick={() => setActiveTab('create')}
                                    className='mt-4'
                                >
                                    <Plus className='h-4 w-4 mr-2' /> Создать сценарий
                                </Button>
                            </div>
                        ) : (
                            <div className='space-y-4'>
                                <div className='grid gap-3'>
                                    {availableScenarios.map((scenario) => (
                                        <div
                                            key={scenario.id}
                                            className={`border p-4 rounded-md cursor-pointer hover:border-primary transition-colors ${
                                                selectedScenarioId === scenario.id
                                                    ? 'border-primary bg-primary/5'
                                                    : ''
                                            }`}
                                            onClick={() => handleSelectScenario(scenario.id!)}
                                        >
                                            <div className='flex justify-between items-center'>
                                                <h3 className='font-medium'>
                                                    Сценарий #{scenario.id}
                                                </h3>
                                                {selectedScenarioId === scenario.id && (
                                                    <span className='text-xs bg-primary text-primary-foreground px-2 py-1 rounded-full'>
                                                        Выбран
                                                    </span>
                                                )}
                                            </div>

                                            {selectedScenarioId === scenario.id &&
                                                renderScenarioDetails(scenario)}
                                        </div>
                                    ))}
                                </div>

                                <DialogFooter>
                                    <Button variant='outline' onClick={() => onOpenChange(false)}>
                                        Отмена
                                    </Button>
                                    <Button
                                        type='button'
                                        onClick={handleApplySelectedScenario}
                                        disabled={!selectedScenarioId}
                                    >
                                        Применить сценарий
                                    </Button>
                                </DialogFooter>
                            </div>
                        )}
                    </TabsContent>

                    <TabsContent value='create' className='space-y-6'>
                        <div className='space-y-4'>
                            <h3 className='text-lg font-medium'>Настройки модели</h3>
                            <div className='space-y-2'>
                                <Label htmlFor='modelName'>Название модели</Label>
                                <Input
                                    id='modelName'
                                    value={scenarioData.model.modelName}
                                    onChange={(e) =>
                                        handleInputChange('model', 'modelName', e.target.value)
                                    }
                                    required
                                />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='systemPrompt'>Системный промпт</Label>
                                <Textarea
                                    id='systemPrompt'
                                    value={scenarioData.model.systemPrompt}
                                    onChange={(e) =>
                                        handleInputChange('model', 'systemPrompt', e.target.value)
                                    }
                                    required
                                />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='temperature'>Температура</Label>
                                <Input
                                    id='temperature'
                                    type='number'
                                    step='0.1'
                                    min='0'
                                    max='1'
                                    value={scenarioData.model.temperature}
                                    onChange={(e) =>
                                        handleInputChange(
                                            'model',
                                            'temperature',
                                            parseFloat(e.target.value)
                                        )
                                    }
                                    required
                                />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='topK'>Top K</Label>
                                <Input
                                    id='topK'
                                    type='number'
                                    min='1'
                                    value={scenarioData.model.topK}
                                    onChange={(e) =>
                                        handleInputChange('model', 'topK', parseInt(e.target.value))
                                    }
                                    required
                                />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='topP'>Top P</Label>
                                <Input
                                    id='topP'
                                    type='number'
                                    step='0.1'
                                    min='0'
                                    max='1'
                                    value={scenarioData.model.topP}
                                    onChange={(e) =>
                                        handleInputChange(
                                            'model',
                                            'topP',
                                            parseFloat(e.target.value)
                                        )
                                    }
                                    required
                                />
                            </div>
                        </div>

                        <div className='space-y-4'>
                            <h3 className='text-lg font-medium'>Векторный поиск</h3>
                            <div className='space-y-2'>
                                <div className='flex items-center space-x-2'>
                                    <Switch
                                        id='searchByQuery'
                                        checked={scenarioData.vectorSearch.searchByQuery}
                                        onCheckedChange={(checked) =>
                                            handleInputChange(
                                                'vectorSearch',
                                                'searchByQuery',
                                                checked
                                            )
                                        }
                                    />
                                    <Label htmlFor='searchByQuery'>Поиск по запросу</Label>
                                </div>
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='threshold'>Порог</Label>
                                <Input
                                    id='threshold'
                                    type='number'
                                    step='0.1'
                                    min='0'
                                    max='1'
                                    value={scenarioData.vectorSearch.threshold}
                                    onChange={(e) =>
                                        handleInputChange(
                                            'vectorSearch',
                                            'threshold',
                                            parseFloat(e.target.value)
                                        )
                                    }
                                    required
                                />
                            </div>
                            <div className='space-y-2'>
                                <Label htmlFor='topN'>Top N</Label>
                                <Input
                                    id='topN'
                                    type='number'
                                    min='1'
                                    value={scenarioData.vectorSearch.topN}
                                    onChange={(e) =>
                                        handleInputChange(
                                            'vectorSearch',
                                            'topN',
                                            parseInt(e.target.value)
                                        )
                                    }
                                    required
                                />
                            </div>
                        </div>

                        <div className='space-y-4'>
                            <h3 className='text-lg font-medium'>Мультизапрос</h3>
                            <div className='space-y-2'>
                                <div className='flex items-center space-x-2'>
                                    <Switch
                                        id='useMultiquery'
                                        checked={useMultiquery}
                                        onCheckedChange={setUseMultiquery}
                                    />
                                    <Label htmlFor='useMultiquery'>Использовать мультизапрос</Label>
                                </div>
                            </div>

                            {useMultiquery && (
                                <>
                                    <div className='space-y-2'>
                                        <Label htmlFor='nQueries'>Количество запросов</Label>
                                        <Input
                                            id='nQueries'
                                            type='number'
                                            min='1'
                                            value={scenarioData.multiQuery.nQueries}
                                            onChange={(e) =>
                                                handleInputChange(
                                                    'multiQuery',
                                                    'nQueries',
                                                    parseInt(e.target.value)
                                                )
                                            }
                                            required
                                        />
                                    </div>
                                    <div className='space-y-2'>
                                        <Label htmlFor='queryModelName'>
                                            Название модели запроса
                                        </Label>
                                        <Input
                                            id='queryModelName'
                                            value={scenarioData.multiQuery.queryModelName}
                                            onChange={(e) =>
                                                handleInputChange(
                                                    'multiQuery',
                                                    'queryModelName',
                                                    e.target.value
                                                )
                                            }
                                            required
                                        />
                                    </div>
                                </>
                            )}
                        </div>

                        <div className='space-y-4'>
                            <h3 className='text-lg font-medium'>Ранжирование</h3>
                            <div className='space-y-2'>
                                <div className='flex items-center space-x-2'>
                                    <Switch
                                        id='useRerank'
                                        checked={useRerank}
                                        onCheckedChange={setUseRerank}
                                    />
                                    <Label htmlFor='useRerank'>Использовать ранжирование</Label>
                                </div>
                            </div>

                            {useRerank && (
                                <>
                                    <div className='space-y-2'>
                                        <Label htmlFor='rerankerModel'>Модель ранжирования</Label>
                                        <Input
                                            id='rerankerModel'
                                            value={scenarioData.reranker.rerankerModel}
                                            onChange={(e) =>
                                                handleInputChange(
                                                    'reranker',
                                                    'rerankerModel',
                                                    e.target.value
                                                )
                                            }
                                            required
                                        />
                                    </div>
                                    <div className='space-y-2'>
                                        <Label htmlFor='rerankerMaxLength'>
                                            Максимальная длина
                                        </Label>
                                        <Input
                                            id='rerankerMaxLength'
                                            type='number'
                                            min='1'
                                            value={scenarioData.reranker.rerankerMaxLength}
                                            onChange={(e) =>
                                                handleInputChange(
                                                    'reranker',
                                                    'rerankerMaxLength',
                                                    parseInt(e.target.value)
                                                )
                                            }
                                            required
                                        />
                                    </div>
                                    <div className='space-y-2'>
                                        <Label htmlFor='rerankerTopK'>Top K</Label>
                                        <Input
                                            id='rerankerTopK'
                                            type='number'
                                            min='1'
                                            value={scenarioData.reranker.topK}
                                            onChange={(e) =>
                                                handleInputChange(
                                                    'reranker',
                                                    'topK',
                                                    parseInt(e.target.value)
                                                )
                                            }
                                            required
                                        />
                                    </div>
                                </>
                            )}
                        </div>

                        <DialogFooter>
                            <Button variant='outline' onClick={() => onOpenChange(false)}>
                                Отмена
                            </Button>
                            <Button
                                type='button'
                                onClick={handleCreateScenario}
                                disabled={isCreating}
                            >
                                {isCreating ? (
                                    <>
                                        <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                                        Создание...
                                    </>
                                ) : (
                                    'Создать сценарий'
                                )}
                            </Button>
                        </DialogFooter>
                    </TabsContent>
                </Tabs>
            </DialogContent>
        </Dialog>
    );
});

export default ScenarioSettingsModal;
