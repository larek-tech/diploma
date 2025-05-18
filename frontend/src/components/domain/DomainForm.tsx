import {DomainApiService} from '@/api/DomainApiService';
import {CreateDomainRequest, Domain, Source} from '@/api/models';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandSeparator,
} from '@/components/ui/command';
import {Drawer, DrawerContent, DrawerFooter, DrawerHeader, DrawerTitle, DrawerTrigger} from '@/components/ui/drawer';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import {Popover, PopoverContent, PopoverTrigger} from '@/components/ui/popover';
import {useToast} from '@/components/ui/use-toast';
import {cn} from '@/lib/utils';
import {Loader2, Plus, X} from 'lucide-react';
import {useRef, useState} from 'react';
import {CreateSourceForm} from './CreateSourceForm';

export const DomainForm = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [isSourceLoading, setIsSourceLoading] = useState(true);
    const [domainTitle, setDomainTitle] = useState('');
    const [selectedSourceIds, setSelectedSourceIds] = useState<number[]>([]);
    const [availableSources, setAvailableSources] = useState<Source[]>([]);
    const [createdDomain, setCreatedDomain] = useState<Domain | null>(null);
    const [error, setError] = useState<string | null>(null);

    const [isDrawerOpen, setIsDrawerOpen] = useState(false);
    const [popoverOpen, setPopoverOpen] = useState(false);
    const { toast } = useToast();

    // Reference to the Command component for focusing
    const commandRef = useRef<HTMLDivElement>(null);

    const loadSources = async () => {
        try {
            setIsSourceLoading(true);
            const response = await DomainApiService.getSources();
            setAvailableSources(response.sources);
        } catch (error) {
            console.error('Failed to load sources:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить источники',
                variant: 'destructive',
            });
        } finally {
            setIsSourceLoading(false);
        }
    };

    // Load sources on component mount
    useState(() => {
        loadSources();
    });

    const handleSourceCreated = (source: Source) => {
        // Add the new source to the available sources list
        setAvailableSources((prev) => [...prev, source]);

        // Automatically select the newly created source
        setSelectedSourceIds((prev) => [...prev, source.id]);

        // Close the source creation drawer
        setIsDrawerOpen(false);

        toast({
            title: 'Успех',
            description: `Источник "${source.title}" успешно создан и выбран`,
            variant: 'default',
        });
    };

    const toggleSource = (sourceId: number) => {
        setSelectedSourceIds((prev) =>
            prev.includes(sourceId) ? prev.filter((id) => id !== sourceId) : [...prev, sourceId]
        );
    };

    const removeSource = (sourceId: number, e?: React.MouseEvent) => {
        if (e) {
            e.stopPropagation();
        }
        setSelectedSourceIds((prev) => prev.filter((id) => id !== sourceId));
    };

    const handleCreateDomain = async () => {
        if (!domainTitle.trim()) {
            setError('Необходимо указать название домена');
            return;
        }

        if (selectedSourceIds.length === 0) {
            setError('Необходимо выбрать хотя бы один источник');
            return;
        }

        try {
            setIsLoading(true);
            setError(null);

            const domainData: CreateDomainRequest = {
                title: domainTitle,
                sourceIds: selectedSourceIds,
            };

            const domain = await DomainApiService.createDomain(domainData);
            setCreatedDomain(domain);

            toast({
                title: 'Успех',
                description: `Домен "${domain.title}" успешно создан`,
                variant: 'default',
            });

            // Reset form after successful creation
            setDomainTitle('');
            setSelectedSourceIds([]);
        } catch (err) {
            setError('Ошибка при создании домена');
            console.error(err);

            toast({
                title: 'Ошибка',
                description: 'Не удалось создать домен',
                variant: 'destructive',
            });
        } finally {
            setIsLoading(false);
        }
    };

    // Get the selected sources
    const selectedSources = availableSources.filter((source) =>
        selectedSourceIds.includes(source.id)
    );

    return (
        <div className='space-y-6'>
            <Card>
                <CardHeader>
                    <CardTitle>Создание нового домена</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className='space-y-6'>
                        <div className='space-y-4'>
                            <div className='space-y-2'>
                                <Label htmlFor='domainTitle'>Название домена</Label>
                                <Input
                                    id='domainTitle'
                                    value={domainTitle}
                                    onChange={(e) => setDomainTitle(e.target.value)}
                                    placeholder='Введите название домена'
                                />
                            </div>

                            {/* Multi-select source picker with create option */}
                            <div className='space-y-2'>
                                <Label>Источники</Label>
                                <div className='flex flex-col gap-2'>
                                    {/* Selected sources badges */}
                                    {selectedSources.length > 0 && (
                                        <div className='flex flex-wrap gap-1'>
                                            {selectedSources.map((source) => (
                                                <Badge
                                                    key={source.id}
                                                    variant='secondary'
                                                    className='flex items-center gap-1'
                                                >
                                                    {source.title}
                                                    <X
                                                        className='h-3 w-3 cursor-pointer'
                                                        onClick={(e) => removeSource(source.id, e)}
                                                    />
                                                </Badge>
                                            ))}
                                        </div>
                                    )}

                                    {/* Source selector */}
                                    <div className='flex gap-2'>
                                        <Popover open={popoverOpen} onOpenChange={setPopoverOpen}>
                                            <PopoverTrigger asChild>
                                                <Button
                                                    variant='outline'
                                                    className='justify-start flex-1 text-left'
                                                >
                                                    {selectedSourceIds.length > 0
                                                        ? `Выбрано: ${selectedSourceIds.length}`
                                                        : 'Выберите источники'}
                                                </Button>
                                            </PopoverTrigger>
                                            <PopoverContent
                                                className='p-0'
                                                align='start'
                                                side='bottom'
                                                sideOffset={10}
                                            >
                                                <Command ref={commandRef}>
                                                    <CommandInput placeholder='Поиск источников...' />
                                                    <CommandList>
                                                        <CommandEmpty>
                                                            {isSourceLoading ? (
                                                                <div className='py-6 text-center text-sm'>
                                                                    <Loader2 className='h-4 w-4 animate-spin mx-auto' />
                                                                    <p className='mt-2'>
                                                                        Загрузка источников...
                                                                    </p>
                                                                </div>
                                                            ) : (
                                                                'Источники не найдены'
                                                            )}
                                                        </CommandEmpty>
                                                        <CommandGroup heading='Доступные источники'>
                                                            {availableSources.map((source) => (
                                                                <CommandItem
                                                                    key={source.id}
                                                                    onSelect={() =>
                                                                        toggleSource(source.id)
                                                                    }
                                                                    className={cn(
                                                                        'flex items-center gap-2',
                                                                        selectedSourceIds.includes(
                                                                            source.id
                                                                        ) && 'bg-primary/10'
                                                                    )}
                                                                >
                                                                    <span
                                                                        className={cn(
                                                                            'h-4 w-4 border rounded-sm',
                                                                            selectedSourceIds.includes(
                                                                                source.id
                                                                            )
                                                                                ? 'bg-primary border-primary'
                                                                                : 'border-input'
                                                                        )}
                                                                    />
                                                                    <span>{source.title}</span>
                                                                </CommandItem>
                                                            ))}
                                                        </CommandGroup>
                                                        <CommandSeparator />
                                                        <CommandGroup>
                                                            <CommandItem
                                                                onSelect={() => {
                                                                    setIsDrawerOpen(true);
                                                                    setPopoverOpen(false);
                                                                }}
                                                                className='text-primary'
                                                            >
                                                                <Plus className='mr-2 h-4 w-4' />
                                                                <span>Создать новый источник</span>
                                                            </CommandItem>
                                                        </CommandGroup>
                                                    </CommandList>
                                                </Command>
                                            </PopoverContent>
                                        </Popover>

                                        <Drawer open={isDrawerOpen} onOpenChange={setIsDrawerOpen}>
                                            <DrawerTrigger asChild>
                                                <Button
                                                    variant='outline'
                                                    className='px-3'
                                                    onClick={() => setIsDrawerOpen(true)}
                                                >
                                                    <Plus className='h-4 w-4' />
                                                </Button>
                                            </DrawerTrigger>
                                            <DrawerContent>
                                                <DrawerHeader>
                                                    <DrawerTitle>
                                                        Создать новый источник
                                                    </DrawerTitle>
                                                </DrawerHeader>
                                                <div className='px-4'>
                                                    <CreateSourceForm
                                                        onSourceCreated={handleSourceCreated}
                                                    />
                                                </div>
                                                <DrawerFooter>
                                                    <Button
                                                        variant='outline'
                                                        onClick={() => setIsDrawerOpen(false)}
                                                    >
                                                        Отмена
                                                    </Button>
                                                </DrawerFooter>
                                            </DrawerContent>
                                        </Drawer>
                                    </div>
                                </div>
                            </div>

                            {error && <div className='text-red-500'>{error}</div>}

                            {createdDomain ? (
                                <div className='flex items-center rounded-md border border-green-200 bg-green-50 p-4 text-sm text-green-800 dark:bg-green-900/30 dark:text-green-200'>
                                    <p>
                                        Домен "{createdDomain.title}" успешно создан (ID:{' '}
                                        {createdDomain.id})
                                    </p>
                                </div>
                            ) : (
                                <Button
                                    onClick={handleCreateDomain}
                                    disabled={isLoading}
                                    className='w-full'
                                >
                                    {isLoading ? (
                                        <>
                                            <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                                            Создание...
                                        </>
                                    ) : (
                                        'Создать домен'
                                    )}
                                </Button>
                            )}
                        </div>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};
