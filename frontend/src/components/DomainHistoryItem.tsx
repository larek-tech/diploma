import { Format } from '@/utils/date-utils';
import { Database } from 'lucide-react';
import { Domain } from '@/api/models';
import { formatDate } from '@/utils/date-utils';

interface DomainHistoryItemProps {
    domain: Domain;
}

const DomainHistoryItem = ({ domain }: DomainHistoryItemProps) => {
    const createdDate = formatDate(
        new Date(domain.createdAt.seconds * 1000),
        Format.DayMonthYearTime
    );

    return (
        <div className='bg-white border border-gray-200 hover:border-primary rounded-lg p-4 transition-colors duration-200 hover:shadow-sm'>
            <div className='flex items-center space-x-3'>
                <div className='flex-shrink-0'>
                    <Database className='h-6 w-6 text-primary' />
                </div>
                <div className='flex-1 min-w-0'>
                    <p className='text-sm font-medium text-gray-900 truncate'>{domain.title}</p>
                    <p className='text-xs text-gray-500'>
                        ID: {domain.id} | Создан: {createdDate}
                    </p>
                    <p className='text-xs text-gray-500'>Источников: {domain.sourceIds.length}</p>
                </div>
            </div>
        </div>
    );
};

export default DomainHistoryItem;
