import {WSMessageType} from '@/api/models';
import Conversation from '@/components/Conversation';
import EmptyChat from '@/components/EmptyChat';
import {Button} from '@/components/ui/button';
import {Skeleton} from '@/components/ui/skeleton';
import {Textarea} from '@/components/ui/textarea';
import {useToast} from '@/components/ui/use-toast';
import {useStores} from '@/hooks/useStores';
import {Pages} from '@/router/constants';
import debounce from 'lodash/debounce';
import {ArrowUpIcon, BookOpen, Database, FilePenIcon, Loader2, MicIcon, Settings, StopCircleIcon} from 'lucide-react';
import {observer} from 'mobx-react-lite';
import {ChangeEvent, KeyboardEvent, useEffect, useRef, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import ScenarioSettingsModal from './ScenarioSettingsModal';

const ChatWrapper = observer(() => {
    const { rootStore } = useStores();
    const { toast } = useToast();
    const { sessionId } = useParams();
    const navigate = useNavigate();
    const [message, setMessage] = useState('');
    const [titleValue, setTitleValue] = useState('');
    const [recognizing, setRecognizing] = useState(false);
    const [isScenarioSettingsOpen, setIsScenarioSettingsOpen] = useState(false);

    const titleInputRef = useRef<HTMLInputElement>(null);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const recognitionRef = useRef<any | null>(null);

    useEffect(() => {
        console.log('sessionId', sessionId);

        if (!sessionId) {
            rootStore
                .createSession()
                .then(() => {
                    rootStore.getSessions();
                })
                .catch(() => {
                    toast({
                        title: 'Ошибка',
                        description: 'Не удалось создать сессию',
                        variant: 'destructive',
                    });
                });
        } else {
            rootStore.getSession({ id: sessionId }).catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось загрузить сессию',
                    variant: 'destructive',
                });

                navigate(`/${Pages.Chat}`, { replace: true });
            });
        }
    }, [rootStore, sessionId, navigate, toast]);

    useEffect(() => {
        if (rootStore.chatError) {
            toast({
                title: 'Ошибка',
                description: rootStore.chatError,
                variant: 'destructive',
            });
        }
    }, [rootStore.chatError, toast]);

    const handleKeyDown = (event: KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            sendMessage();
        }
    };

    const sendMessage = () => {
        if (message.trim() && !rootStore.isChatDisabled && rootStore.websocket?.readyState === 1) {
            rootStore.sendMessage({
                type: WSMessageType.Query,
                content: message.trim(),
                domainID: rootStore.selectedDomainId ?? undefined,
                scenarioID: rootStore.selectedScenarioId ?? undefined,
            });
            setMessage('');
        } else if (!rootStore.selectedDomainId) {
            toast({
                title: 'Внимание',
                description: 'Выберите домен знаний для чата в боковой панели',
                variant: 'destructive',
            });
        } else if (rootStore.hasScenarios() && !rootStore.selectedScenarioId) {
            toast({
                title: 'Внимание',
                description: 'Выберите сценарий для выбранного домена знаний',
                variant: 'destructive',
            });
        }
    };

    const debouncedRenameSession = debounce((sessionId: string, title: string) => {
        rootStore
            .renameSession({ id: sessionId, title })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Сессия успешно переименована',
                    variant: 'default',
                });
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось переименовать сессию',
                    variant: 'destructive',
                });
            });
    }, 1000);

    const onTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
        setTitleValue(event.target.value);

        const title = event.target?.value;
        if (rootStore.activeSessionId) {
            debouncedRenameSession(rootStore.activeSessionId, title);
        }
    };

    const startRecognition = () => {
        if (!('SpeechRecognition' in window) && !('webkitSpeechRecognition' in window)) {
            toast({
                title: 'Ошибка',
                description: 'Ваш браузер не поддерживает Web Speech API',
                variant: 'destructive',
            });
            return;
        }

        if (!recognitionRef.current) {
            const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;
            const recognition = new SpeechRecognition();
            recognition.continuous = false;
            recognition.interimResults = true;
            recognition.lang = 'ru-RU';

            recognition.onstart = () => {
                setRecognizing(true);
            };

            // eslint-disable-next-line
            // @ts-ignore
            recognition.onresult = (event) => {
                let interimTranscript = '';
                for (let i = event.resultIndex; i < event.results.length; ++i) {
                    if (event.results[i].isFinal) {
                        setMessage(event.results[i][0].transcript);
                    } else {
                        interimTranscript += event.results[i][0].transcript;
                    }
                }
                setMessage((prev) => prev + interimTranscript);
            };

            // eslint-disable-next-line
            // @ts-ignore
            recognition.onerror = (event) => {
                console.error(event.error);
                toast({
                    title: 'Ошибка',
                    description: 'Ошибка распознавания речи',
                    variant: 'destructive',
                });
            };

            recognition.onend = () => {
                setRecognizing(false);
            };

            recognitionRef.current = recognition;
        }

        recognitionRef.current.start();
    };

    const stopRecognition = () => {
        recognitionRef.current?.stop();
    };

    return (
        <>
            <div className='chat'>
                <div className='flex items-center justify-between'>
                    <div className='flex items-center gap-2 group'>
                        <input
                            ref={titleInputRef}
                            type='text'
                            className='bg-transparent text-lg font-medium focus:outline-none'
                            value={titleValue || rootStore.activeSession?.title || 'Новый чат'}
                            onChange={(event) => onTitleChange(event)}
                        />
                        <Button
                            variant='ghost'
                            size='icon'
                            className='rounded-full hover:bg-gray-100 dark:hover:bg-[#1e293b] transition-colors'
                            onClick={() => {
                                titleInputRef.current?.focus();
                            }}
                        >
                            <FilePenIcon className='w-5 h-5 text-gray-500 dark:text-gray-400 group-hover:text-gray-700 dark:group-hover:text-gray-300' />
                        </Button>
                    </div>

                    <div className='flex items-center gap-4 text-sm'>
                        <div className='flex items-center gap-1'>
                            <Database className='h-4 w-4 text-primary' />
                            <span
                                className={
                                    rootStore.selectedDomainId
                                        ? 'text-primary-600 font-medium'
                                        : 'text-gray-500'
                                }
                            >
                                {rootStore.getSelectedDomainTitle()}
                            </span>
                        </div>

                        {rootStore.selectedDomainId && rootStore.hasScenarios() && (
                            <div className='flex items-center gap-1'>
                                <span
                                    className={
                                        rootStore.selectedScenarioId
                                            ? 'text-primary-600 font-medium'
                                            : 'text-gray-500'
                                    }
                                >
                                    <div className='flex items-center gap-1'>
                                        <BookOpen className='h-4 w-4 text-primary' />
                                        <span
                                            className={
                                                rootStore.selectedScenarioId
                                                    ? 'text-primary-600 font-medium'
                                                    : 'text-gray-500'
                                            }
                                        >
                                            {rootStore.selectedScenarioId
                                                ? `Сценарий #${rootStore.getSelectedScenarioId()}`
                                                : 'Выберите сценарий'}
                                        </span>
                                        <Button
                                            variant='ghost'
                                            size='icon'
                                            className='rounded-full hover:bg-gray-100 dark:hover:bg-[#1e293b] transition-colors p-1'
                                            onClick={() => setIsScenarioSettingsOpen(true)}
                                        >
                                            <Settings className='w-4 h-4 text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300' />
                                        </Button>
                                    </div>
                                </span>
                            </div>
                        )}
                    </div>
                </div>

                <div className='chat__messages-area flex flex-col'>
                    <div className='max-w-5xl w-full flex-1 mx-auto flex flex-col items-start gap-8 px-4 py-8'>
                        {rootStore.activeSessionLoading ? (
                            Array.from({ length: 3 }).map((_, i) => (
                                <div
                                    key={i}
                                    className='flex items-start gap-4 animate-pulse w-full'
                                >
                                    <Skeleton className='bg-slate-200 w-8 h-8 rounded-full' />
                                    <div className='grid gap-1 flex-1'>
                                        <Skeleton className='bg-slate-200 h-8 w-24' />
                                        <Skeleton
                                            className={`bg-slate-200 w-full ${
                                                Math.random() > 0.5 ? 'h-8' : 'h-12'
                                            }`}
                                        />
                                    </div>
                                </div>
                            ))
                        ) : rootStore.activeDisplayedSession?.messages.length ? (
                            rootStore.activeDisplayedSession?.messages.map((conversation, i) => (
                                <Conversation
                                    key={i}
                                    conversation={conversation}
                                    isLastConversation={
                                        i ===
                                        (rootStore.activeDisplayedSession?.messages.length || 0) - 1
                                    }
                                />
                            ))
                        ) : (
                            <EmptyChat />
                        )}
                    </div>

                    <div className='max-w-5xl w-full sticky bottom-0 mx-auto py-4 flex flex-col gap-2 px-4 dark:bg-[#0f172a] bg-neutral-100'>
                        {!rootStore.selectedDomainId && (
                            <div className='text-amber-600 text-xs mb-1 flex items-center'>
                                <span className='mr-1'>⚠️</span> Выберите домен знаний в боковой
                                панели для начала общения
                            </div>
                        )}

                        {rootStore.selectedDomainId &&
                            rootStore.hasScenarios() &&
                            !rootStore.selectedScenarioId && (
                                <div className='text-amber-600 text-xs mb-1 flex items-center'>
                                    <span className='mr-1'>⚠️</span> Выберите сценарий для
                                    выбранного домена
                                </div>
                            )}

                        <div className='relative'>
                            <Textarea
                                onChange={(e) => setMessage(e.target.value)}
                                onKeyDown={(event) => handleKeyDown(event)}
                                value={message}
                                placeholder='Напишите в чат...'
                                name='message'
                                id='message'
                                rows={1}
                                className='min-h-[48px] rounded-2xl resize-none p-4 border border-gray-300 shadow-sm pr-16 dark:border-gray-800'
                            />

                            <Button
                                type='submit'
                                size='icon'
                                variant='outline'
                                className={`absolute top-3 right-12 w-8 h-8 ${
                                    !recognizing && 'hidden'
                                }`}
                                onClick={stopRecognition}
                            >
                                <StopCircleIcon className='w-4 h-4' />
                                <span className='sr-only'>Остановить запрос</span>
                            </Button>

                            <Button
                                type='button'
                                size='icon'
                                className={`absolute top-3 right-12 w-8 h-8 ${
                                    recognizing && 'hidden'
                                }`}
                                onClick={startRecognition}
                                disabled={recognizing}
                            >
                                <MicIcon className='w-4 h-4' />
                                <span className='sr-only'>Начать распознавание</span>
                            </Button>

                            <Button
                                type='submit'
                                size='icon'
                                className='absolute top-3 right-3 w-8 h-8'
                                onClick={sendMessage}
                                disabled={rootStore.isChatDisabled || !rootStore.websocket}
                            >
                                {rootStore.isModelAnswering ? (
                                    <Loader2 className='absolute h-4 w-4 animate-spin' />
                                ) : (
                                    <ArrowUpIcon className='w-4 h-4' />
                                )}

                                <span className='sr-only'>Отправить</span>
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            <ScenarioSettingsModal
                open={isScenarioSettingsOpen}
                onOpenChange={setIsScenarioSettingsOpen}
            />
        </>
    );
});

export default ChatWrapper;
