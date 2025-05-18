import {Domain} from '@/api/models';
import {useStores} from '@/hooks/useStores';
import {cn} from '@/lib/utils';
import {Format, formatDate} from '@/utils/date-utils';
import {CheckCircle2, Database} from 'lucide-react';
import {observer} from 'mobx-react-lite';

interface DomainHistoryItemProps {
    domain: Domain;
}

const DomainHistoryItem = observer(({ domain }: DomainHistoryItemProps) => {
    const { rootStore } = useStores();
    const isSelected = rootStore.selectedDomainId === domain.id;

    const createdDate = formatDate(
        new Date(domain.createdAt.seconds * 1000),
        Format.DayMonthYearTime
    );

    // const handleDomainClick = () => {
    //     rootStore.setSelectedDomain(domain.id);
    // };

    return (
        <div
            className={cn(
                'bg-white border rounded-lg p-4 transition-colors duration-200 hover:shadow-sm cursor-pointer relative',
                isSelected ? 'border-primary bg-primary/5' : 'border-gray-200 hover:border-primary'
            )}
            // onClick={handleDomainClick}
        >
            <div className='flex items-center space-x-3'>
                <div className='flex-shrink-0'>
                    <Database className='h-6 w-6 text-primary' />
                </div>
                <div className='flex-1 min-w-0'>
                    <div className='flex items-center gap-2'>
                        <p className='text-sm font-medium text-gray-900 truncate'>{domain.title}</p>
                        {isSelected && <CheckCircle2 className='h-4 w-4 text-primary' />}
                    </div>
                    <p className='text-xs text-gray-500'>
                        ID: {domain.id} <br /> Создан: {createdDate}
                    </p>
                    <div className='flex justify-between items-center'>
                        <p className='text-xs text-gray-500'>
                            Источников: {domain.sourceIds.length} <br /> Сценариев:{' '}
                            {domain?.scenarioIds?.length}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
});

export default DomainHistoryItem;
