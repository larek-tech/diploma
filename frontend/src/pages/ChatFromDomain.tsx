import ChatWrapper from '@/components/ChatWrapper';
import {useStores} from '@/hooks/useStores';
import {useEffect} from 'react';
import {useParams} from 'react-router-dom';

const ChatFromDomain = () => {
    const { domainId } = useParams();
    const { rootStore } = useStores();

    useEffect(() => {
        if (domainId) {
            rootStore.setSelectedDomain(+domainId);
        }
    }, [domainId, rootStore]);

    return <ChatWrapper />;
};

export default ChatFromDomain;
