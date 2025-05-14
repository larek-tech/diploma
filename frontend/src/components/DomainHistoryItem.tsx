import { Format } from '@/utils/date-utils';
import { Database, CheckCircle2, ChevronDown } from 'lucide-react';
import { Domain } from '@/api/models';
import { formatDate } from '@/utils/date-utils';
import { useStores } from '@/hooks/useStores';
import { observer } from 'mobx-react-lite';
import { cn } from '@/lib/utils';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';

interface DomainHistoryItemProps {
    domain: Domain;
}

const DomainHistoryItem = observer(({ domain }: DomainHistoryItemProps) => {
    const { rootStore } = useStores();
    const isSelected = rootStore.selectedDomainId === domain.id;
    const [dropdownOpen, setDropdownOpen] = useState(false);

    const createdDate = formatDate(
        new Date(domain.createdAt.seconds * 1000),
        Format.DayMonthYearTime
    );

    const handleDomainClick = () => {
        rootStore.setSelectedDomain(domain.id);
    };

    const handleScenarioSelect = (scenarioId: number) => {
        rootStore.setSelectedScenario(scenarioId);
        setDropdownOpen(false);
    };

    return (
        <div
            className={cn(
                'bg-white border rounded-lg p-4 transition-colors duration-200 hover:shadow-sm cursor-pointer relative',
                isSelected ? 'border-primary bg-primary/5' : 'border-gray-200 hover:border-primary'
            )}
            onClick={handleDomainClick}
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
                        ID: {domain.id} | Создан: {createdDate}
                    </p>
                    <div className='flex justify-between items-center'>
                        <p className='text-xs text-gray-500'>
                            Источников: {domain.sourceIds.length} | Сценариев:{' '}
                            {domain?.scenarioIds?.length}
                        </p>

                        {isSelected && domain.scenarioIds.length > 0 && (
                            <div onClick={(e) => e.stopPropagation()} className='mt-1'>
                                <DropdownMenu open={dropdownOpen} onOpenChange={setDropdownOpen}>
                                    <DropdownMenuTrigger asChild>
                                        <Button
                                            variant='outline'
                                            size='sm'
                                            className='h-7 text-xs px-2 py-0'
                                        >
                                            <span className='mr-1'>
                                                {rootStore.selectedScenarioId
                                                    ? `Сценарий #${rootStore.selectedScenarioId}`
                                                    : 'Выбрать сценарий'}
                                            </span>
                                            <ChevronDown className='h-3 w-3' />
                                        </Button>
                                    </DropdownMenuTrigger>
                                    <DropdownMenuContent align='end' className='w-48'>
                                        {domain.scenarioIds.map((scenarioId) => (
                                            <DropdownMenuItem
                                                key={scenarioId}
                                                className={cn(
                                                    'text-xs cursor-pointer',
                                                    rootStore.selectedScenarioId === scenarioId &&
                                                        'bg-primary/10'
                                                )}
                                                onClick={() => handleScenarioSelect(scenarioId)}
                                            >
                                                <span className='flex items-center gap-2'>
                                                    {rootStore.selectedScenarioId ===
                                                        scenarioId && (
                                                        <CheckCircle2 className='h-3 w-3 text-primary' />
                                                    )}
                                                    Сценарий #{scenarioId}
                                                </span>
                                            </DropdownMenuItem>
                                        ))}
                                    </DropdownMenuContent>
                                </DropdownMenu>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
});

export default DomainHistoryItem;
