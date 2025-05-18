import {useStores} from '@/hooks/useStores';
import {Pages} from '@/router/constants';
import {Database} from 'lucide-react';
import {observer} from 'mobx-react-lite';
import {useEffect} from 'react';
import {Link} from 'react-router-dom';
import DomainHistoryItem from './DomainHistoryItem';
import {Button} from './ui/button';
import {Skeleton} from './ui/skeleton';

const DomainsHistory = observer(() => {
    const { rootStore } = useStores();

    useEffect(() => {
        // Загружаем домены при первом рендере
        if (rootStore.domains.length === 0 && !rootStore.domainsLoading) {
            rootStore.getDomains();
        }
    }, [rootStore]);

    const handleLoadMore = () => {
        if (!rootStore.domainsLoading && rootStore.hasMoreDomains) {
            rootStore.getDomains();
        }
    };

    return (
        <div className='w-full mx-auto'>
            <div className='flex items-center justify-between mb-2'>
                <h3 className='text-lg font-medium text-gray-700'>Домены знаний</h3>
            </div>
            <div className='space-y-2'>
                {rootStore.domainsLoading && rootStore.domains.length === 0 ? (
                    // Показываем скелетон загрузки только при первой загрузке
                    Array.from({ length: 3 }).map((_, i) => (
                        <div
                            key={i}
                            className='bg-gray-100 rounded-lg p-4 transition-colors duration-300'
                        >
                            <div className='flex items-center space-x-3'>
                                <Skeleton className='h-6 w-6 rounded-full' />
                                <div className='space-y-2 flex-1'>
                                    <Skeleton className='h-4 w-32' />
                                    <Skeleton className='h-3 w-24' />
                                </div>
                            </div>
                        </div>
                    ))
                ) : rootStore.domains.length > 0 ? (
                    <>
                        {rootStore.domains.map((domain) => (
                            <Link
                                key={domain.id}
                                to={`/${Pages.Domain}/${domain.id}`}
                                className='block cursor-pointer'
                            >
                                <DomainHistoryItem key={domain.id} domain={domain} />
                            </Link>
                        ))}

                        {rootStore.hasMoreDomains && (
                            <Button
                                variant='outline'
                                className='w-full mt-2'
                                onClick={handleLoadMore}
                                disabled={rootStore.domainsLoading}
                            >
                                {rootStore.domainsLoading ? 'Загрузка...' : 'Загрузить еще'}
                            </Button>
                        )}
                    </>
                ) : (
                    <div className='text-center py-4 text-gray-500'>
                        <Database className='mx-auto h-8 w-8 mb-2 text-gray-400' />
                        <p>Нет доступных доменов</p>
                        <p className='text-sm'>Создайте новый домен знаний</p>
                    </div>
                )}
            </div>
        </div>
    );
});

export default DomainsHistory;
