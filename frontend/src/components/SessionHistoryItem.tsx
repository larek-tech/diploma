import {ShortChatSession} from '@/api/models';
import {useStores} from '@/hooks/useStores';
import {Pages} from '@/router/constants';
import {Format, formatDate} from '@/utils/date-utils';
import {Loader2, TrashIcon} from 'lucide-react';
import {observer} from 'mobx-react-lite';
import {useState} from 'react';
import {Link, useNavigate} from 'react-router-dom';
import {useToast} from './ui/use-toast';

type Props = {
    session: ShortChatSession;
};

const SessionHistoryItem = observer(({ session }: Props) => {
    const { rootStore } = useStores();
    const [isDeleting, setIsDeleting] = useState(false);
    const { toast } = useToast();
    const navigate = useNavigate();

    const onDeleteSession = (id: string) => {
        setIsDeleting(true);

        rootStore
            .deleteSession({ id })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Сессия успешно удалена',
                    variant: 'default',
                });

                rootStore.getSessions();

                navigate(`/${Pages.Chat}`);
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось удалить сессию',
                    variant: 'destructive',
                });
            })
            .finally(() => {
                setIsDeleting(false);
            });
    };

    return (
        <Link to={`/chat/${session.id}`}>
            <div
                className={`session-history-item hover:bg-gray-200 rounded-lg p-4 transition-colors duration-300 hover:cursor-pointer ${
                    rootStore.activeSessionId === session.id ? 'bg-slate-200' : 'bg-gray-100'
                }`}
            >
                <h3 className='text-sm font-medium mb-1'>
                    {session.title || session.id.slice(0, 8)}...
                </h3>
                <div className='flex items-center justify-between'>
                    <p className='text-gray-500 text-sm'>
                        {formatDate(
                            new Date(session.createdAt.seconds * 1000),
                            Format.DayMonthYearTime
                        )}
                    </p>
                    <button
                        onClick={(event) => {
                            event.stopPropagation();
                            event.preventDefault();

                            onDeleteSession(session.id);
                        }}
                        disabled={isDeleting}
                        className='session-history-item__delete-button text-gray-400 hover:text-gray-600 focus:outline-none group hidden'
                    >
                        {isDeleting ? (
                            <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                        ) : (
                            <TrashIcon className='h-4 w-4 group-hover:block' />
                        )}
                        <span className='sr-only'>Удалить из истории</span>
                    </button>
                </div>
            </div>
        </Link>
    );
});

export default SessionHistoryItem;
