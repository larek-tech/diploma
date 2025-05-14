import { ClipboardIcon } from 'lucide-react';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Button } from './ui/button';
import { useToast } from './ui/use-toast';
import MarkdownPreview from '@uiw/react-markdown-preview';

type ModelMessageProps = {
    incomingMessage: string;
    isLastMessage: boolean;
};

const ModelMessage = ({ incomingMessage }: ModelMessageProps) => {
    const { toast } = useToast();

    console.log(incomingMessage);

    const getModelResonse = () => {
        return (
            <>
                <div className='flex w-full flex-col gap-5'>
                    <div className='prose prose-stone overflow-x-scroll markdown'>
                        <div>
                            <MarkdownPreview source={incomingMessage} style={{ padding: 16 }} />
                        </div>
                    </div>
                </div>
                <div className='flex items-center gap-2 py-2'>
                    <Button
                        variant='ghost'
                        size='icon'
                        className='w-4 h-4 hover:bg-transparent text-stone-400 hover:text-stone-900'
                        onClick={() => {
                            navigator.clipboard.writeText(incomingMessage);
                            toast({
                                title: 'Скопировано',
                                description: 'Текст ответа скопирован в буфер обмена',
                            });
                        }}
                    >
                        <ClipboardIcon className='w-4 h-4' />
                        <span className='sr-only'>Копировать</span>
                    </Button>
                </div>
            </>
        );
    };

    return (
        <div className='flex items-start gap-4 w-full'>
            <Avatar className='border w-8 h-8'>
                <AvatarFallback>MT</AvatarFallback>
            </Avatar>
            <div className='gap-1 mt-2 w-full'>
                <div className='font-bold'>Ответ модели</div>

                {getModelResonse()}
            </div>
        </div>
    );
};

export default ModelMessage;
