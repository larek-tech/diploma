import MarkdownPreview from '@uiw/react-markdown-preview';
import {ClipboardIcon, ExternalLinkIcon} from 'lucide-react';
import {useEffect, useState} from 'react';
import {Avatar, AvatarFallback} from './ui/avatar';
import {Badge} from './ui/badge';
import {Button} from './ui/button';
import {Card, CardContent} from './ui/card';
import {useToast} from './ui/use-toast';

type ModelMessageProps = {
    incomingMessage: string;
    isLastMessage: boolean;
};

type Source =
    | {
          metadata?: {
              resourceUrl?: string;
              source?: string;
              title?: string;
              name?: string;
              description?: string;
              url?: string;
          };
          id?: string;
          url?: string; // Alternative field name that might be used
          name?: string; // Alternative field name that might be used
          title?: string; // Alternative field name for title
          source?: string; // Alternative field name for source URL
          resourceUrl?: string; // Alternative field name for resource URL
      }
    | string; // Allow source to be just a string (URL)

const ModelMessage = ({ incomingMessage }: ModelMessageProps) => {
    const { toast } = useToast();
    const [messageContent, setMessageContent] = useState<string>('');
    const [sources, setSources] = useState<Source[]>([]);

    useEffect(() => {
        if (!incomingMessage) {
            setMessageContent('');
            setSources([]);
            return;
        }

        console.log('Processing incoming message:', incomingMessage);

        // Check for different variations of "Источники:" string
        const sourceKeywords = ['Источники:', 'источники:', 'Sources:', 'sources:'];
        let contentPart = incomingMessage;
        let sourcesPart = '';

        // Find the first occurrence of any source keyword
        for (const keyword of sourceKeywords) {
            const splitIndex = incomingMessage.indexOf(keyword);
            if (splitIndex !== -1) {
                contentPart = incomingMessage.substring(0, splitIndex).trim();
                sourcesPart = incomingMessage.substring(splitIndex + keyword.length).trim();
                console.log(`Found sources with keyword "${keyword}" at position ${splitIndex}`);
                console.log('Sources part:', sourcesPart);
                break;
            }
        }

        // Set the content part of the message
        setMessageContent(contentPart);

        // If we found sources text, try to parse it
        if (sourcesPart) {
            try {
                // Try to extract and parse JSON from the sources part
                // Look for something that looks like an array or object
                const jsonMatch = sourcesPart.match(/(\[.*\]|\{.*\})/s);

                if (jsonMatch && jsonMatch[0]) {
                    console.log('Found JSON in sources:', jsonMatch[0]);

                    // Try to parse the JSON
                    try {
                        const sourcesArray = JSON.parse(jsonMatch[0]);
                        console.log('Parsed sources:', sourcesArray);

                        // Handle both array and object formats
                        if (Array.isArray(sourcesArray)) {
                            console.log('Setting sources array:', sourcesArray);
                            setSources(sourcesArray);
                        } else if (typeof sourcesArray === 'object') {
                            // If it's a single object, wrap it in an array
                            console.log('Setting single source object:', sourcesArray);
                            setSources([sourcesArray]);
                        }
                    } catch (jsonError) {
                        console.error('JSON parse error:', jsonError);

                        // Try alternative parsing approaches
                        try {
                            // Try to fix common JSON issues
                            const fixedJson = jsonMatch[0]
                                .replace(/'/g, '"') // Replace single quotes with double quotes
                                .replace(/([{,])\s*([a-zA-Z0-9_]+)\s*:/g, '$1"$2":') // Add quotes to keys
                                .replace(/:\s*([a-zA-Z0-9_]+)\s*([,}])/g, ':"$1"$2'); // Add quotes to string values

                            console.log('Attempting to parse fixed JSON:', fixedJson);
                            const sourcesArray = JSON.parse(fixedJson);

                            if (Array.isArray(sourcesArray)) {
                                setSources(sourcesArray);
                            } else if (typeof sourcesArray === 'object') {
                                setSources([sourcesArray]);
                            }
                        } catch (fixError) {
                            console.error('Failed to parse even with fixes:', fixError);
                            // If JSON parsing fails, include the sources text as part of the message
                            setMessageContent(contentPart + '\n\nИсточники:\n' + sourcesPart);
                        }
                    }
                } else {
                    console.log('No valid JSON pattern found in sources part');

                    // Try to extract URLs directly from the text
                    const urlRegex = /(https?:\/\/[^\s]+)/g;
                    const urls = sourcesPart.match(urlRegex);

                    if (urls && urls.length > 0) {
                        console.log('Found URLs directly in text:', urls);
                        const extractedSources = urls.map((url, index) => ({
                            metadata: { resourceUrl: url },
                            id: `extracted-${index}`,
                        }));
                        setSources(extractedSources);
                    } else {
                        // If no valid JSON or URLs found, still display the sources text as content
                        setMessageContent(contentPart + '\n\nИсточники:\n' + sourcesPart);
                    }
                }
            } catch (e) {
                console.error('Failed to parse sources JSON:', e);
                // If JSON parsing fails, include the sources text as part of the message
                setMessageContent(contentPart + '\n\nИсточники:\n' + sourcesPart);
            }
        }
    }, [incomingMessage]);

    // Helper function to get URL from source object regardless of structure
    const getSourceUrl = (source: Source): string | null => {
        // Check for direct string (some APIs might just return array of URLs)
        if (typeof source === 'string') return source;

        // Check top-level properties first (new format: {"title": "gitflic", "resourceUrl": "https://..."})
        if (source.resourceUrl) return source.resourceUrl;
        if (source.url) return source.url;
        if (source.source) return source.source;

        // Check various nested properties
        if (source.metadata?.resourceUrl) return source.metadata.resourceUrl;
        if (source.metadata?.source) return source.metadata.source;
        if (source.metadata?.url) return source.metadata.url;

        // Check if the entire source might be a URL string
        const sourceStr = String(source);
        if (sourceStr.match(/^https?:\/\//)) return sourceStr;

        return null;
    };

    // Helper function to get title/name from source object
    const getSourceTitle = (source: Source, index: number): string => {
        // If source is just a string (URL), extract domain
        if (typeof source === 'string') {
            try {
                const url = new URL(source);
                return url.hostname;
            } catch {
                return `Источник ${index + 1}`;
            }
        }

        // Check top-level properties first (new format: {"title": "gitflic", "resourceUrl": "https://..."})
        if (source.title) return source.title;
        if (source.name) return source.name;

        // Check various nested properties for title/name
        if (source.metadata?.title) return source.metadata.title;
        if (source.metadata?.name) return source.metadata.name;
        if (source.id) return `Источник ${source.id}`;

        return `Источник ${index + 1}`;
    };

    const getModelResonse = () => {
        return (
            <>
                <div className='flex w-full flex-col gap-5'>
                    <div className='prose prose-stone overflow-x-scroll markdown'>
                        <div>
                            <MarkdownPreview source={messageContent} style={{ padding: 16 }} />
                        </div>
                    </div>

                    {sources.length > 0 && (
                        <Card className='mt-4'>
                            <CardContent className='pt-4'>
                                <h4 className='font-medium mb-2'>Источники:</h4>
                                <div className='flex flex-wrap gap-2'>
                                    {sources.map((source, index) => {
                                        const url = getSourceUrl(source);
                                        const title = getSourceTitle(source, index);

                                        console.log(`Source ${index}:`, { source, url, title });

                                        return (
                                            <Badge
                                                key={index}
                                                variant='secondary'
                                                className='flex items-center gap-1 p-2 my-1'
                                            >
                                                {url ? (
                                                    <a
                                                        href={url}
                                                        target='_blank'
                                                        rel='noopener noreferrer'
                                                        className='text-white hover:underline flex items-center'
                                                    >
                                                        {title}
                                                        <ExternalLinkIcon className='ml-1 w-3 h-3' />
                                                    </a>
                                                ) : (
                                                    <span>{title}</span>
                                                )}

                                                {typeof source !== 'string' &&
                                                    source.metadata?.description && (
                                                        <span className='ml-2 text-xs text-gray-500 truncate max-w-[200px]'>
                                                            {source.metadata.description}
                                                        </span>
                                                    )}
                                            </Badge>
                                        );
                                    })}
                                </div>
                            </CardContent>
                        </Card>
                    )}
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
