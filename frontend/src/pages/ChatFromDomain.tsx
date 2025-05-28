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

    useEffect(() => {
        if (domainId) {
            // Сбрасываем состояние готовности при смене домена
            setAllSourcesReady(false);
            rootStore.setSelectedDomain(+domainId);
        }
    }, [domainId, rootStore]);

    const handleStatusCheck = (allReady: boolean) => {
        setAllSourcesReady(allReady);
    };

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
