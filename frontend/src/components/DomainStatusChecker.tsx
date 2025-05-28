import {DomainApiService} from '@/api/DomainApiService';
import {Source, SourceStatus} from '@/api/models';
import {Badge} from '@/components/ui/badge';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {Skeleton} from '@/components/ui/skeleton';
import {useStores} from '@/hooks/useStores';
import {AlertCircle, CheckCircle, Loader2, XCircle} from 'lucide-react';
import {observer} from 'mobx-react-lite';
import {useEffect, useState} from 'react';

interface DomainStatusCheckerProps {
    onStatusCheck?: (allReady: boolean) => void;
}

const getStatusIcon = (status: SourceStatus) => {
    switch (status) {
        case SourceStatus.Ready:
            return <CheckCircle className='w-4 h-4 text-green-500' />;
        case SourceStatus.Parsing:
            return <Loader2 className='w-4 h-4 text-yellow-500 animate-spin' />;
        case SourceStatus.Failed:
            return <XCircle className='w-4 h-4 text-red-500' />;
        case SourceStatus.Undefind:
            return <AlertCircle className='w-4 h-4 text-gray-500' />;
        default:
            return <AlertCircle className='w-4 h-4 text-gray-500' />;
    }
};

const getStatusText = (status: SourceStatus) => {
    switch (status) {
        case SourceStatus.Ready:
            return 'Готов';
        case SourceStatus.Parsing:
            return 'Обработка';
        case SourceStatus.Failed:
            return 'Ошибка';
        case SourceStatus.Undefind:
            return 'Не определен';
        default:
            return 'Неизвестно';
    }
};

const getStatusVariant = (status: SourceStatus) => {
    switch (status) {
        case SourceStatus.Ready:
            return 'default';
        case SourceStatus.Parsing:
            return 'secondary';
        case SourceStatus.Failed:
            return 'destructive';
        case SourceStatus.Undefind:
            return 'outline';
        default:
            return 'outline';
    }
};

const DomainStatusChecker = observer(({ onStatusCheck }: DomainStatusCheckerProps) => {
    const { rootStore } = useStores();
    const [sources, setSources] = useState<Source[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const checkDomainSources = async () => {
            if (!rootStore.selectedDomain) {
                setSources([]);
                onStatusCheck?.(false); // Явно передаем false когда нет домена
                return;
            }

            setLoading(true);
            try {
                // Если у домена нет источников, считаем его готовым
                if (
                    !rootStore.selectedDomain.sourceIds ||
                    rootStore.selectedDomain.sourceIds.length === 0
                ) {
                    setSources([]);
                    onStatusCheck?.(true);
                    return;
                }

                // Получаем информацию о каждом источнике в домене
                const sourcePromises = rootStore.selectedDomain.sourceIds.map((sourceId) =>
                    DomainApiService.getSource(sourceId)
                );

                const fetchedSources = await Promise.all(sourcePromises);
                setSources(fetchedSources);

                // Проверяем, все ли источники готовы
                const allReady = fetchedSources.every(
                    (source) => source.status === SourceStatus.Ready
                );
                onStatusCheck?.(allReady);
            } catch (error) {
                console.error('Ошибка при получении статуса источников:', error);
                onStatusCheck?.(false);
            } finally {
                setLoading(false);
            }
        };

        checkDomainSources();
    }, [rootStore.selectedDomain, onStatusCheck]);

    // Сброс состояния при смене домена
    useEffect(() => {
        setSources([]);
        setLoading(false);
    }, [rootStore.selectedDomain?.id]);

    // Периодическая проверка статуса источников, если есть источники в процессе обработки
    useEffect(() => {
        const hasProcessingSources = sources.some(
            (source) =>
                source.status === SourceStatus.Parsing || source.status === SourceStatus.Undefind
        );

        if (!hasProcessingSources) {
            return;
        }

        const interval = setInterval(async () => {
            if (!rootStore.selectedDomain) {
                return;
            }

            try {
                // Если у домена нет источников, считаем его готовым
                if (
                    !rootStore.selectedDomain.sourceIds ||
                    rootStore.selectedDomain.sourceIds.length === 0
                ) {
                    onStatusCheck?.(true);
                    return;
                }

                const sourcePromises = rootStore.selectedDomain.sourceIds.map((sourceId) =>
                    DomainApiService.getSource(sourceId)
                );

                const fetchedSources = await Promise.all(sourcePromises);
                setSources(fetchedSources);

                const allReady = fetchedSources.every(
                    (source) => source.status === SourceStatus.Ready
                );
                onStatusCheck?.(allReady);
            } catch (error) {
                console.error('Ошибка при обновлении статуса источников:', error);
            }
        }, 5000); // Проверяем каждые 5 секунд

        return () => clearInterval(interval);
    }, [sources, rootStore.selectedDomain, onStatusCheck]);

    if (!rootStore.selectedDomain) {
        return null;
    }

    const allReady = sources.every((source) => source.status === SourceStatus.Ready);
    const hasFailedSources = sources.some((source) => source.status === SourceStatus.Failed);

    return (
        <Card className='w-full max-w-4xl mx-auto'>
            <CardHeader>
                <div className='flex items-center justify-between'>
                    <CardTitle className='flex items-center gap-2'>
                        <span>Статус домена: {rootStore.selectedDomain.title}</span>
                        {loading && (
                            <div className='flex items-center gap-1'>
                                <Loader2 className='w-5 h-5 animate-spin text-blue-500' />
                                <span className='text-sm text-blue-500 animate-pulse'>
                                    Проверка...
                                </span>
                            </div>
                        )}
                    </CardTitle>
                    {allReady && (
                        <Badge variant='default' className='bg-green-100 text-green-800'>
                            <CheckCircle className='w-3 h-3 mr-1' />
                            Готов к использованию
                        </Badge>
                    )}
                    {hasFailedSources && (
                        <Badge variant='destructive'>
                            <XCircle className='w-3 h-3 mr-1' />
                            Есть ошибки
                        </Badge>
                    )}
                    {!allReady && !hasFailedSources && (
                        <Badge variant='secondary' className='animate-pulse'>
                            <Loader2 className='w-3 h-3 mr-1 animate-spin' />
                            Обработка...
                        </Badge>
                    )}
                </div>
            </CardHeader>
            <CardContent>
                <div className='space-y-3'>
                    <h4 className='text-sm font-medium text-gray-700 dark:text-gray-300'>
                        Источники данных:
                    </h4>
                    <div className='space-y-2'>
                        {loading && sources.length === 0 ? (
                            // Скелетоны во время загрузки
                            Array.from({ length: 3 }).map((_, index) => (
                                <div
                                    key={index}
                                    className='flex items-center justify-between p-3 border rounded-lg bg-gray-50 dark:bg-gray-800 animate-pulse'
                                >
                                    <div className='flex items-center gap-3'>
                                        <Skeleton className='w-4 h-4 rounded-full' />
                                        <div className='space-y-1'>
                                            <Skeleton className='h-4 w-32' />
                                            <Skeleton className='h-3 w-20' />
                                        </div>
                                    </div>
                                    <Skeleton className='h-6 w-16 rounded-full' />
                                </div>
                            ))
                        ) : sources.length === 0 ? (
                            // Сообщение когда нет источников
                            <div className='p-4 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg'>
                                <div className='flex items-center gap-2 text-blue-800 dark:text-blue-200'>
                                    <CheckCircle className='w-4 h-4' />
                                    <p className='text-sm font-medium'>
                                        У домена нет источников данных
                                    </p>
                                </div>
                                <p className='text-xs text-blue-700 dark:text-blue-300 mt-1'>
                                    Домен готов к использованию.
                                </p>
                            </div>
                        ) : (
                            sources.map((source) => (
                                <div
                                    key={source.id}
                                    className={`flex items-center justify-between p-3 border rounded-lg bg-gray-50 dark:bg-gray-800 transition-all duration-300 ${
                                        source.status === SourceStatus.Parsing
                                            ? 'ring-2 ring-yellow-200 dark:ring-yellow-800'
                                            : ''
                                    }`}
                                >
                                    <div className='flex items-center gap-3'>
                                        {getStatusIcon(source.status)}
                                        <div>
                                            <p className='font-medium text-sm'>{source.title}</p>
                                            <p className='text-xs text-gray-500'>
                                                ID: {source.id} • Тип: {source.typ}
                                            </p>
                                        </div>
                                    </div>
                                    <Badge
                                        variant={getStatusVariant(source.status)}
                                        className={
                                            source.status === SourceStatus.Parsing
                                                ? 'animate-pulse'
                                                : ''
                                        }
                                    >
                                        {getStatusText(source.status)}
                                    </Badge>
                                </div>
                            ))
                        )}
                    </div>

                    {!allReady && sources.length > 0 && (
                        <div className='mt-4 p-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg animate-pulse'>
                            <div className='flex items-center gap-2 text-amber-800 dark:text-amber-200'>
                                <AlertCircle className='w-4 h-4 animate-pulse' />
                                <p className='text-sm font-medium'>
                                    Домен не готов к использованию
                                </p>
                            </div>
                            <p className='text-xs text-amber-700 dark:text-amber-300 mt-1'>
                                Дождитесь завершения обработки всех источников перед началом работы
                                с чатом.
                            </p>
                        </div>
                    )}
                </div>
            </CardContent>
        </Card>
    );
});

export default DomainStatusChecker;
