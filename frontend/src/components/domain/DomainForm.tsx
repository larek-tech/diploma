import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { DomainApiService } from '@/api/DomainApiService';
import { CreateDomainRequest, Domain, Source } from '@/api/models';
import { SourceSelector } from './SourceSelector';
import { CreateSourceForm } from './CreateSourceForm';
import { ScenarioForm } from './ScenarioForm';
import { Loader2 } from 'lucide-react';

enum CreateDomainStep {
    SOURCES = 'sources',
    SCENARIO = 'scenario',
}

export const DomainForm = () => {
    const [currentStep, setCurrentStep] = useState<CreateDomainStep>(CreateDomainStep.SOURCES);
    const [isLoading, setIsLoading] = useState(false);
    const [domainTitle, setDomainTitle] = useState('');
    const [selectedSourceIds, setSelectedSourceIds] = useState<number[]>([]);
    const [newSources, setNewSources] = useState<Source[]>([]);
    const [createdDomain, setCreatedDomain] = useState<Domain | null>(null);
    const [error, setError] = useState<string | null>(null);

    const handleSourcesSelected = (sourceIds: number[]) => {
        setSelectedSourceIds(sourceIds);
    };

    const handleSourceCreated = (source: Source) => {
        setNewSources((prev) => [...prev, source]);
        setSelectedSourceIds((prev) => [...prev, source.id]);
    };

    const handleCreateDomain = async () => {
        if (!domainTitle.trim()) {
            setError('Необходимо указать название домена');
            return;
        }

        if (selectedSourceIds.length === 0) {
            setError('Необходимо выбрать хотя бы один источник');
            return;
        }

        try {
            setIsLoading(true);
            setError(null);

            const domainData: CreateDomainRequest = {
                title: domainTitle,
                sourceIds: selectedSourceIds,
            };

            const domain = await DomainApiService.createDomain(domainData);
            setCreatedDomain(domain);
            setCurrentStep(CreateDomainStep.SCENARIO);
        } catch (err) {
            setError('Ошибка при создании домена');
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    };

    const handleScenarioCreated = () => {
        // Можно добавить обработку успешного создания сценария
        // например, редирект на страницу домена или уведомление
    };

    const handleTabChange = (value: string) => {
        setCurrentStep(value as CreateDomainStep);
    };

    return (
        <div className='space-y-6'>
            <Card>
                <CardHeader>
                    <CardTitle>Создание нового домена</CardTitle>
                </CardHeader>
                <CardContent>
                    <Tabs defaultValue={currentStep} onValueChange={handleTabChange}>
                        <TabsList className='grid w-full grid-cols-2'>
                            <TabsTrigger
                                value={CreateDomainStep.SOURCES}
                                disabled={
                                    currentStep === CreateDomainStep.SCENARIO && !createdDomain
                                }
                            >
                                Источники
                            </TabsTrigger>
                            <TabsTrigger
                                value={CreateDomainStep.SCENARIO}
                                disabled={!createdDomain}
                            >
                                Сценарий
                            </TabsTrigger>
                        </TabsList>

                        <TabsContent value={CreateDomainStep.SOURCES} className='space-y-6'>
                            <div className='space-y-4'>
                                <div className='space-y-2'>
                                    <Label htmlFor='domainTitle'>Название домена</Label>
                                    <Input
                                        id='domainTitle'
                                        value={domainTitle}
                                        onChange={(e) => setDomainTitle(e.target.value)}
                                        placeholder='Введите название домена'
                                    />
                                </div>

                                {/* Выбор существующих источников */}
                                <SourceSelector onSourcesSelected={handleSourcesSelected} />

                                {/* Список новых созданных источников */}
                                {newSources.length > 0 && (
                                    <div className='space-y-2'>
                                        <h3 className='text-lg font-medium'>Созданные источники</h3>
                                        <div className='space-y-1'>
                                            {newSources.map((source) => (
                                                <div key={source.id} className='text-sm'>
                                                    • {source.title} (ID: {source.id})
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Форма создания нового источника */}
                                <div className='pt-4'>
                                    <h3 className='text-lg font-medium mb-4'>
                                        Создать новый источник
                                    </h3>
                                    <CreateSourceForm onSourceCreated={handleSourceCreated} />
                                </div>

                                {error && <div className='text-red-500'>{error}</div>}

                                <Button onClick={handleCreateDomain} disabled={isLoading}>
                                    {isLoading ? (
                                        <>
                                            <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                                            Создание...
                                        </>
                                    ) : (
                                        'Создать домен и перейти к сценарию'
                                    )}
                                </Button>
                            </div>
                        </TabsContent>

                        <TabsContent value={CreateDomainStep.SCENARIO}>
                            {createdDomain ? (
                                <>
                                    <div className='mb-4'>
                                        <p>
                                            Домен "{createdDomain.title}" успешно создан (ID:{' '}
                                            {createdDomain.id})
                                        </p>
                                        <p>Теперь настройте сценарий для вашего домена:</p>
                                    </div>
                                    <ScenarioForm
                                        domainId={createdDomain.id}
                                        onScenarioCreated={handleScenarioCreated}
                                    />
                                </>
                            ) : (
                                <div>Сначала создайте домен</div>
                            )}
                        </TabsContent>
                    </Tabs>
                </CardContent>
            </Card>
        </div>
    );
};
