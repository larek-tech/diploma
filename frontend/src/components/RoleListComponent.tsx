import {RoleForEntity} from '@/api/models/roles';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from '@/components/ui/dropdown-menu';
import {Skeleton} from '@/components/ui/skeleton';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from '@/components/ui/table';
import {Format, formatDate} from '@/utils/date-utils';
import {MoreVertical, RefreshCw, Shield, Trash2} from 'lucide-react';
import React, {useState} from 'react';

interface RoleListComponentProps {
    roles: RoleForEntity[];
    loading: boolean;
    onDelete: (roleId: number) => Promise<void>;
    onRefresh: () => Promise<void>;
}

const RoleListComponent: React.FC<RoleListComponentProps> = ({
    roles,
    loading,
    onDelete,
    onRefresh,
}) => {
    const [deletingRoleId, setDeletingRoleId] = useState<number | null>(null);

    const handleDeleteRole = async (role: RoleForEntity) => {
        if (!confirm(`Вы уверены, что хотите удалить роль "${role.name}"?`)) {
            return;
        }

        try {
            setDeletingRoleId(role.id);
            await onDelete(role.id);
        } catch (error) {
            console.error('Error deleting role:', error);
        } finally {
            setDeletingRoleId(null);
        }
    };

    if (loading) {
        return (
            <div className='space-y-4'>
                <div className='flex justify-between items-center'>
                    <Skeleton className='h-8 w-48' />
                    <Skeleton className='h-10 w-32' />
                </div>
                <div className='space-y-2'>
                    {Array.from({ length: 5 }).map((_, index) => (
                        <div key={index} className='flex items-center space-x-4'>
                            <Skeleton className='h-12 w-12 rounded-full' />
                            <div className='space-y-2 flex-1'>
                                <Skeleton className='h-4 w-48' />
                                <Skeleton className='h-3 w-32' />
                            </div>
                            <Skeleton className='h-10 w-10' />
                        </div>
                    ))}
                </div>
            </div>
        );
    }

    return (
        <div className='space-y-4'>
            <div className='flex justify-between items-center'>
                <div className='flex items-center gap-2'>
                    <Shield className='h-5 w-5' />
                    <span className='font-medium'>Всего ролей: {roles.length}</span>
                </div>
                <Button variant='outline' size='sm' onClick={onRefresh} disabled={loading}>
                    <RefreshCw className='mr-2 h-4 w-4' />
                    Обновить
                </Button>
            </div>

            {roles.length === 0 ? (
                <div className='text-center py-12'>
                    <Shield className='mx-auto h-12 w-12 text-muted-foreground' />
                    <h3 className='mt-4 text-lg font-semibold'>Нет ролей</h3>
                    <p className='mt-2 text-muted-foreground'>
                        Создайте первую роль для начала работы с системой разрешений.
                    </p>
                </div>
            ) : (
                <div className='border rounded-lg'>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Название роли</TableHead>
                                <TableHead>ID</TableHead>
                                <TableHead>Дата создания</TableHead>
                                <TableHead>Статус</TableHead>
                                <TableHead className='w-[100px]'>Действия</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {roles.map((role) => (
                                <TableRow key={role.id}>
                                    <TableCell className='font-medium'>
                                        <div className='flex items-center gap-2'>
                                            <Shield className='h-4 w-4 text-primary' />
                                            {role.name}
                                        </div>
                                    </TableCell>
                                    <TableCell className='text-muted-foreground'>
                                        #{role.id}
                                    </TableCell>
                                    <TableCell>
                                        {formatDate(
                                            new Date(role.createdAt.seconds * 1000),
                                            Format.DayMonthYearTime
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant='secondary'>Активна</Badge>
                                    </TableCell>
                                    <TableCell>
                                        <DropdownMenu>
                                            <DropdownMenuTrigger asChild>
                                                <Button
                                                    variant='ghost'
                                                    size='sm'
                                                    className='h-8 w-8 p-0'
                                                >
                                                    <MoreVertical className='h-4 w-4' />
                                                    <span className='sr-only'>Открыть меню</span>
                                                </Button>
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent align='end'>
                                                <DropdownMenuItem
                                                    onClick={() => handleDeleteRole(role)}
                                                    disabled={deletingRoleId === role.id}
                                                    className='text-red-600'
                                                >
                                                    <Trash2 className='mr-2 h-4 w-4' />
                                                    {deletingRoleId === role.id
                                                        ? 'Удаление...'
                                                        : 'Удалить'}
                                                </DropdownMenuItem>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </div>
            )}
        </div>
    );
};

export default RoleListComponent;
