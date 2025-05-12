import { useState, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { DomainApiService } from '@/api/DomainApiService';
import { CreateSourceRequest, Source, SourceType } from '@/api/models';
import { Loader2, Upload, File } from 'lucide-react';
import { useDropzone } from 'react-dropzone';

interface CreateSourceFormProps {
    onSourceCreated: (source: Source) => void;
}

const sourceSchema = z.object({
    title: z.string().min(1, 'Название источника обязательно'),
    typ: z.number().min(1).max(4),
    content: z.string().min(1, 'Содержимое обязательно'),
    credentials: z.string().optional(),
    updateParams: z
        .object({
            everyPeriod: z.number().optional(),
            cron: z
                .object({
                    minute: z.number().optional(),
                    hour: z.number().optional(),
                    dayOfMonth: z.number().optional(),
                    month: z.number().optional(),
                    dayOfWeek: z.number().optional(),
                })
                .optional(),
        })
        .optional(),
});

type SourceFormValues = z.infer<typeof sourceSchema>;

export const CreateSourceForm = ({ onSourceCreated }: CreateSourceFormProps) => {
    const [isLoading, setIsLoading] = useState(false);
    const [showCronFields, setShowCronFields] = useState(false);
    const [showCredentialsField, setShowCredentialsField] = useState(false);
    const [showFileUpload, setShowFileUpload] = useState(false);
    const [selectedFiles, setSelectedFiles] = useState<File[] | null>(null);

    const {
        register,
        handleSubmit,
        setValue,
        reset,
        watch,
        formState: { errors },
    } = useForm<SourceFormValues>({
        resolver: zodResolver(sourceSchema),
        defaultValues: {
            typ: SourceType.TypeWeb,
        },
    });

    const handleTypeChange = (value: string) => {
        const typeValue = parseInt(value, 10);
        setValue('typ', typeValue);

        // Показываем поле для credentials только если тип TypeWithCredentials
        setShowCredentialsField(typeValue === SourceType.TypeWithCredentials);

        // Показываем загрузку файла для TypeSingleFile и TypeArchivedFiles
        setShowFileUpload(
            typeValue === SourceType.TypeSingleFile || typeValue === SourceType.TypeArchivedFiles
        );

        // Сбрасываем значение content при смене типа
        if (showFileUpload) {
            setValue('content', '');
            setSelectedFiles(null);
        }
    };

    const handleCronToggle = (usesCron: boolean) => {
        setShowCronFields(usesCron);

        if (!usesCron) {
            // Сбрасываем значения cron полей
            setValue('updateParams', { everyPeriod: 0 });
        } else {
            // Устанавливаем значения по умолчанию для cron
            setValue('updateParams', {
                cron: {
                    minute: 0,
                    hour: 0,
                    dayOfMonth: 0,
                    month: 0,
                    dayOfWeek: 0,
                },
            });
        }
    };

    // Функция кодирования файла в base64
    const encodeFileToBase64 = (file: File): Promise<string> => {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.readAsDataURL(file);
            reader.onload = () => {
                if (typeof reader.result === 'string') {
                    // Удаляем префикс "data:application/..." из строки base64
                    const base64String = reader.result.split(',')[1];
                    resolve(base64String);
                } else {
                    reject(new Error('Не удалось прочитать файл как строку'));
                }
            };
            reader.onerror = (error) => reject(error);
        });
    };

    // Обработчик изменения файла
    const handleFileChange = useCallback(
        async (files: File[] | null) => {
            setSelectedFiles(files);

            if (files && files.length > 0) {
                try {
                    const base64 = await encodeFileToBase64(files[0]);
                    setValue('content', base64);
                } catch (error) {
                    console.error('Ошибка кодирования файла в base64:', error);
                }
            }
        },
        [setValue]
    );

    // Пользовательский компонент для загрузки файла
    const FileUploadArea = () => {
        const onDrop = useCallback((acceptedFiles: File[]) => {
            handleFileChange(acceptedFiles);
        }, []);

        const { getRootProps, getInputProps, isDragActive } = useDropzone({
            onDrop,
            maxFiles: 1,
            multiple: false,
            accept: {
                'application/zip': ['.zip'],
                'application/pdf': ['.pdf'],
                'application/msword': ['.doc'],
                'application/vnd.openxmlformats-officedocument.wordprocessingml.document': [
                    '.docx',
                ],
                'text/plain': ['.txt'],
                'application/json': ['.json'],
                'text/html': ['.html', '.htm'],
                'text/csv': ['.csv'],
                'application/xml': ['.xml'],
                'application/vnd.ms-excel': ['.xls'],
                'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet': ['.xlsx'],
            },
        });

        const sourceType = watch('typ');
        const acceptArchiveOnly = sourceType === SourceType.TypeArchivedFiles;

        return (
            <div
                {...getRootProps()}
                className={`border-2 border-dashed rounded-md p-6 text-center cursor-pointer hover:bg-gray-50 ${
                    isDragActive ? 'border-primary' : 'border-gray-300'
                }`}
            >
                <input {...getInputProps()} />
                <div className='flex flex-col items-center justify-center'>
                    <Upload className='w-10 h-10 mb-2 text-gray-400' />
                    {selectedFiles && selectedFiles.length > 0 ? (
                        <div className='flex items-center space-x-2'>
                            <File className='w-5 h-5' />
                            <span>{selectedFiles[0].name}</span>
                        </div>
                    ) : (
                        <>
                            <p className='text-sm text-gray-500'>
                                {acceptArchiveOnly
                                    ? 'Перетащите архив сюда или нажмите для выбора'
                                    : 'Перетащите файл сюда или нажмите для выбора'}
                            </p>
                            <p className='text-xs text-gray-400 mt-1'>
                                {acceptArchiveOnly
                                    ? 'Только ZIP файлы'
                                    : 'Поддерживаются PDF, DOC, DOCX, TXT, JSON, HTML, CSV, XML, XLS, XLSX'}
                            </p>
                        </>
                    )}
                </div>
            </div>
        );
    };

    const onSubmit = async (data: SourceFormValues) => {
        try {
            setIsLoading(true);

            // Для типов, требующих текстовый ввод (Web и S3)
            if (data.typ === SourceType.TypeWeb || data.typ === SourceType.TypeWithCredentials) {
                // Кодируем content в base64, если это еще не сделано
                if (!data.content.startsWith('data:') && data.content.trim() !== '') {
                    data.content = btoa(data.content);
                }
            }

            // Кодируем credentials в base64, если заполнено
            if (
                data.credentials &&
                !data.credentials.startsWith('data:') &&
                data.credentials.trim() !== ''
            ) {
                data.credentials = btoa(data.credentials);
            }

            const source = await DomainApiService.createSource(data as CreateSourceRequest);
            onSourceCreated(source);
            reset();
            setSelectedFiles(null);
        } catch (error) {
            console.error('Ошибка при создании источника:', error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className='space-y-4'>
            <div className='space-y-2'>
                <Label htmlFor='title'>Название источника</Label>
                <Input id='title' {...register('title')} />
                {errors.title && <p className='text-sm text-red-500'>{errors.title.message}</p>}
            </div>

            <div className='space-y-2'>
                <Label htmlFor='type'>Тип источника</Label>
                <Select
                    onValueChange={handleTypeChange}
                    defaultValue={SourceType.TypeWeb.toString()}
                >
                    <SelectTrigger>
                        <SelectValue placeholder='Выберите тип источника' />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectItem value={SourceType.TypeWeb.toString()}>Веб-сайт</SelectItem>
                        <SelectItem value={SourceType.TypeSingleFile.toString()}>
                            Один файл
                        </SelectItem>
                        <SelectItem value={SourceType.TypeArchivedFiles.toString()}>
                            Архив с файлами
                        </SelectItem>
                        <SelectItem value={SourceType.TypeWithCredentials.toString()}>
                            S3 с учетными данными
                        </SelectItem>
                    </SelectContent>
                </Select>
            </div>

            {showFileUpload ? (
                <div className='space-y-2'>
                    <Label>Загрузите файл</Label>
                    <FileUploadArea />
                    <input type='hidden' {...register('content')} />
                    {errors.content && (
                        <p className='text-sm text-red-500'>{errors.content.message}</p>
                    )}
                </div>
            ) : (
                <div className='space-y-2'>
                    <Label htmlFor='content'>
                        {watch('typ') === SourceType.TypeWeb ? 'URL веб-сайта' : 'Содержимое'}
                    </Label>
                    <Textarea
                        id='content'
                        {...register('content')}
                        placeholder={
                            watch('typ') === SourceType.TypeWeb
                                ? 'Введите URL веб-сайта'
                                : showCredentialsField
                                ? 'Ссылка на S3 хранилище'
                                : 'Содержимое источника'
                        }
                    />
                    {errors.content && (
                        <p className='text-sm text-red-500'>{errors.content.message}</p>
                    )}
                </div>
            )}

            {showCredentialsField && (
                <div className='space-y-2'>
                    <Label htmlFor='credentials'>Учетные данные</Label>
                    <Textarea
                        id='credentials'
                        {...register('credentials')}
                        placeholder='Секретный ключ для доступа к S3'
                    />
                </div>
            )}

            <div className='space-y-2'>
                <Label>Настройки обновления</Label>
                <div className='flex items-center space-x-2'>
                    <Button
                        type='button'
                        variant={showCronFields ? 'default' : 'outline'}
                        onClick={() => handleCronToggle(true)}
                    >
                        Настроить CRON
                    </Button>
                    <Button
                        type='button'
                        variant={!showCronFields ? 'default' : 'outline'}
                        onClick={() => handleCronToggle(false)}
                    >
                        Период обновления
                    </Button>
                </div>
            </div>

            {showCronFields ? (
                <div className='grid grid-cols-3 gap-2'>
                    <div>
                        <Label htmlFor='minute'>Минута</Label>
                        <Input
                            id='minute'
                            type='number'
                            min={0}
                            max={59}
                            {...register('updateParams.cron.minute', { valueAsNumber: true })}
                        />
                    </div>
                    <div>
                        <Label htmlFor='hour'>Час</Label>
                        <Input
                            id='hour'
                            type='number'
                            min={0}
                            max={23}
                            {...register('updateParams.cron.hour', { valueAsNumber: true })}
                        />
                    </div>
                    <div>
                        <Label htmlFor='dayOfMonth'>День месяца</Label>
                        <Input
                            id='dayOfMonth'
                            type='number'
                            min={1}
                            max={31}
                            {...register('updateParams.cron.dayOfMonth', { valueAsNumber: true })}
                        />
                    </div>
                    <div>
                        <Label htmlFor='month'>Месяц</Label>
                        <Input
                            id='month'
                            type='number'
                            min={1}
                            max={12}
                            {...register('updateParams.cron.month', { valueAsNumber: true })}
                        />
                    </div>
                    <div>
                        <Label htmlFor='dayOfWeek'>День недели</Label>
                        <Input
                            id='dayOfWeek'
                            type='number'
                            min={0}
                            max={6}
                            {...register('updateParams.cron.dayOfWeek', { valueAsNumber: true })}
                        />
                    </div>
                </div>
            ) : (
                <div>
                    <Label htmlFor='everyPeriod'>Период (в секундах)</Label>
                    <Input
                        id='everyPeriod'
                        type='number'
                        min={0}
                        {...register('updateParams.everyPeriod', { valueAsNumber: true })}
                    />
                </div>
            )}

            <Button type='submit' disabled={isLoading} className='w-full'>
                {isLoading ? (
                    <>
                        <Loader2 className='mr-2 h-4 w-4 animate-spin' />
                        Создание...
                    </>
                ) : (
                    'Создать источник'
                )}
            </Button>
        </form>
    );
};
