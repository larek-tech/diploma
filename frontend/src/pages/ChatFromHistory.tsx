import ChatWrapper from '@/components/ChatWrapper';
import {useToast} from '@/components/ui/use-toast';
import {useStores} from '@/hooks/useStores';
import {useEffect} from 'react';
import {useParams} from 'react-router-dom';

const ChatFromHistory = () => {
    const { sessionId } = useParams();
    const { rootStore } = useStores();
    const { toast } = useToast();

    useEffect(() => {
        if (sessionId) {
            rootStore
                .getSession({ id: sessionId })
                .then(() => {
                    const session = rootStore.activeSession;

                    const domainId = session?.content[0].query.domainId;

                    if (domainId) {
                        rootStore.setSelectedDomain(domainId);
                    } else {
                        toast({
                            title: 'Ошибка',
                            description: 'Не удалось загрузить домен',
                            variant: 'destructive',
                        });
                    }
                })
                .catch(() => {
                    toast({
                        title: 'Ошибка',
                        description: 'Не удалось загрузить сессию',
                        variant: 'destructive',
                    });
                });
        }
    }, [sessionId, rootStore, toast]);

    return <ChatWrapper />;
};

export default ChatFromHistory;
