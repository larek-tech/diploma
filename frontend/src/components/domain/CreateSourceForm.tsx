import { useState } from 'react';
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
import { Loader2 } from 'lucide-react';

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

    const {
        register,
        handleSubmit,
        setValue,
        reset,
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

    const onSubmit = async (data: SourceFormValues) => {
        try {
            setIsLoading(true);

            // Конвертируем content в base64 если это еще не сделано
            if (!data.content.startsWith('data:')) {
                data.content = btoa(data.content);
            }

            // Конвертируем credentials в base64 если заполнено
            if (data.credentials && !data.credentials.startsWith('data:')) {
                data.credentials = btoa(data.credentials);
            }

            const source = await DomainApiService.createSource(data as CreateSourceRequest);
            onSourceCreated(source);
            reset();
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

            <div className='space-y-2'>
                <Label htmlFor='content'>Содержимое</Label>
                <Textarea
                    id='content'
                    {...register('content')}
                    placeholder={
                        showCredentialsField
                            ? 'Ссылка на S3 хранилище (будет закодирована в base64)'
                            : 'Ссылка или содержимое (будет закодировано в base64)'
                    }
                />
                {errors.content && <p className='text-sm text-red-500'>{errors.content.message}</p>}
            </div>

            {showCredentialsField && (
                <div className='space-y-2'>
                    <Label htmlFor='credentials'>Учетные данные</Label>
                    <Textarea
                        id='credentials'
                        {...register('credentials')}
                        placeholder='Секретный ключ для доступа к S3 (будет закодирован в base64)'
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
