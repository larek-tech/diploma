import { useEffect, useState } from 'react';
import { Checkbox } from '@/components/ui/checkbox';
import { DomainApiService } from '@/api/DomainApiService';
import { Source } from '@/api/models';
import { Loader2 } from 'lucide-react';

interface SourceSelectorProps {
    onSourcesSelected: (selectedSources: number[]) => void;
}

export const SourceSelector = ({ onSourcesSelected }: SourceSelectorProps) => {
    const [sources, setSources] = useState<Source[]>([]);
    const [selectedSourceIds, setSelectedSourceIds] = useState<number[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchSources = async () => {
            try {
                setLoading(true);
                const response = await DomainApiService.getSources();
                setSources(response.sources);
            } catch (err) {
                setError('Не удалось загрузить источники');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchSources();
    }, []);

    const handleSourceChange = (sourceId: number, checked: boolean) => {
        setSelectedSourceIds((prev) => {
            if (checked) {
                return [...prev, sourceId];
            } else {
                return prev.filter((id) => id !== sourceId);
            }
        });
    };

    useEffect(() => {
        onSourcesSelected(selectedSourceIds);
    }, [selectedSourceIds, onSourcesSelected]);

    if (loading) {
        return (
            <div className='flex justify-center items-center h-40'>
                <Loader2 className='h-8 w-8 animate-spin' />
            </div>
        );
    }

    if (error) {
        return <div className='text-red-500'>{error}</div>;
    }

    if (sources.length === 0) {
        return <div>Источников не найдено. Создайте новый источник.</div>;
    }

    return (
        <div className='space-y-4'>
            <h3 className='text-lg font-medium'>Выберите существующие источники</h3>
            <div className='space-y-2'>
                {sources.map((source) => (
                    <div key={source.id} className='flex items-center space-x-2'>
                        <Checkbox
                            id={`source-${source.id}`}
                            checked={selectedSourceIds.includes(source.id)}
                            onCheckedChange={(checked) =>
                                handleSourceChange(source.id, checked === true)
                            }
                        />
                        <label
                            htmlFor={`source-${source.id}`}
                            className='text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70'
                        >
                            {source.title} (ID: {source.id})
                        </label>
                    </div>
                ))}
            </div>
        </div>
    );
};
