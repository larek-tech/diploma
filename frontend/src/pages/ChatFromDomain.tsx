import ChatWrapper from '@/components/ChatWrapper';
import DomainStatusChecker from '@/components/DomainStatusChecker';
import {useStores} from '@/hooks/useStores';
import {Loader2} from 'lucide-react';
import {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';

const ChatFromDomain = () => {
    const { domainId } = useParams();
    const { rootStore } = useStores();
    const [allSourcesReady, setAllSourcesReady] = useState(false);
    const [isCheckingDomain, setIsCheckingDomain] = useState(true);

    useEffect(() => {
        const loadDomain = async () => {
            if (domainId) {
                console.log('Loading domain:', domainId);
                // Сбрасываем состояние готовности при смене домена
                setAllSourcesReady(false);
                setIsCheckingDomain(true);

                try {
                    await rootStore.setSelectedDomain(+domainId);
                    console.log('Domain loaded:', rootStore.selectedDomain);
                } catch (error) {
                    console.error('Ошибка при загрузке домена:', error);
                } finally {
                    setIsCheckingDomain(false);
                    console.log('Finished loading domain, isCheckingDomain set to false');
                }
            }
        };

        loadDomain();
    }, [domainId, rootStore]);

    const handleStatusCheck = (allReady: boolean) => {
        setAllSourcesReady(allReady);

        // Подключаем WebSocket когда все источники готовы
        if (allReady && rootStore.activeSessionId && !rootStore.websocket) {
            rootStore.connectWebSocketIfReady();
        }
    };

    // Показываем лоадер только если домен еще не загружен
    if (isCheckingDomain) {
        console.log('Showing loader because isCheckingDomain is true');
        return (
            <div className='container mx-auto py-8 flex justify-center items-center'>
                <div className='flex items-center gap-2'>
                    <Loader2 className='w-6 h-6 animate-spin text-blue-500' />
                    <span className='text-lg text-gray-600'>Загрузка домена...</span>
                </div>
            </div>
        );
    }

    if (!rootStore.selectedDomain) {
        console.log('Showing loader because selectedDomain is not set');
        return (
            <div className='container mx-auto py-8 flex justify-center items-center'>
                <div className='flex items-center gap-2'>
                    <Loader2 className='w-6 h-6 animate-spin text-blue-500' />
                    <span className='text-lg text-gray-600'>Домен не найден...</span>
                </div>
            </div>
        );
    }

    // Если не все источники готовы, показываем только статус
    if (!allSourcesReady) {
        return (
            <div className='container mx-auto py-8'>
                <DomainStatusChecker onStatusCheck={handleStatusCheck} />
            </div>
        );
    }

    // Если все источники готовы, показываем чат
    return <ChatWrapper />;
};

export default ChatFromDomain;
