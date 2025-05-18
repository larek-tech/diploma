import {DomainApiService} from '@/api/DomainApiService';
import {CreateDomainRequest, Source} from '@/api/models';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import {Pages} from '@/router/constants';
import {Loader2} from 'lucide-react';
import {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {CreateSourceForm} from './CreateSourceForm';
import {SourceSelector} from './SourceSelector';

export const DomainForm = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [domainTitle, setDomainTitle] = useState('');
    const [selectedSourceIds, setSelectedSourceIds] = useState<number[]>([]);
    const [newSources, setNewSources] = useState<Source[]>([]);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

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
            const domainId = domain.id;

            navigate(`/${Pages.Domain}/${domainId}`);
        } catch (err) {
            setError('Ошибка при создании домена');
            console.error(err);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className='space-y-6'>
            <Card>
                <CardHeader>
                    <CardTitle>Создание нового домена</CardTitle>
                </CardHeader>
                <CardContent>
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
                            <h3 className='text-lg font-medium mb-4'>Создать новый источник</h3>
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
                                'Создать домен'
                            )}
                        </Button>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
