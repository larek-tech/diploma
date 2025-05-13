import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { DomainApiService } from '@/api/DomainApiService';
import { Scenario } from '@/api/models';
import { Loader2 } from 'lucide-react';

interface ScenarioFormProps {
    domainId: number;
    onScenarioCreated: (scenario: unknown) => void;
}

export const ScenarioForm = ({ domainId, onScenarioCreated }: ScenarioFormProps) => {
    const [isLoading, setIsLoading] = useState(false);
    const [useMultiquery, setUseMultiquery] = useState(false);
    const [useRerank, setUseRerank] = useState(false);
    const [scenarioData, setScenarioData] = useState<Scenario>({
        domainId,
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

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            setIsLoading(true);

            // Обновляем значения useMultiquery и useRerank в сценарии
            const updatedScenario = {
                ...scenarioData,
                multiQuery: {
                    ...scenarioData.multiQuery,
                    useMultiquery,
                },
                reranker: {
                    ...scenarioData.reranker,
                    useRerank,
                },
            };

            const result = await DomainApiService.createScenario(updatedScenario);
            onScenarioCreated(result);
        } catch (error) {
            console.error('Ошибка при создании сценария:', error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className='space-y-6'>
            <div className='space-y-4'>
                <h3 className='text-lg font-medium'>Настройки модели</h3>
                <div className='space-y-2'>
                    <Label htmlFor='modelName'>Название модели</Label>
                    <Input
                        id='modelName'
                        value={scenarioData.model.modelName}
                        onChange={(e) => handleInputChange('model', 'modelName', e.target.value)}
                        required
                    />
                </div>
                <div className='space-y-2'>
                    <Label htmlFor='systemPrompt'>Системный промпт</Label>
                    <Textarea
                        id='systemPrompt'
                        value={scenarioData.model.systemPrompt}
                        onChange={(e) => handleInputChange('model', 'systemPrompt', e.target.value)}
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
                            handleInputChange('model', 'temperature', parseFloat(e.target.value))
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
                            handleInputChange('model', 'topP', parseFloat(e.target.value))
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
                                handleInputChange('vectorSearch', 'searchByQuery', checked)
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
                            handleInputChange('vectorSearch', 'topN', parseInt(e.target.value))
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
                            <Label htmlFor='queryModelName'>Название модели запроса</Label>
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
                        <Switch id='useRerank' checked={useRerank} onCheckedChange={setUseRerank} />
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
                                    handleInputChange('reranker', 'rerankerModel', e.target.value)
                                }
                                required
                            />
                        </div>
                        <div className='space-y-2'>
                            <Label htmlFor='rerankerMaxLength'>Максимальная длина</Label>
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
                                    handleInputChange('reranker', 'topK', parseInt(e.target.value))
                                }
                                required
                            />
                        </div>
                    </>
                )}
            </div>

            <Button type='submit' disabled={isLoading} className='w-full'>
                {isLoading ? (
                    <>
                        <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                        Создание...
                    </>
                ) : (
                    'Создать сценарий'
                )}
            </Button>
        </form>
    );
};
