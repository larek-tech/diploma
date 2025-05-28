import {PermittedRoles, PermittedUsers} from '@/api/models';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card';
import {DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger} from '@/components/ui/dropdown-menu';
import {Skeleton} from '@/components/ui/skeleton';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from '@/components/ui/table';
import {useToast} from '@/components/ui/use-toast';
import {Edit, Plus, Trash2, User, Users} from 'lucide-react';
import React, {useEffect, useState} from 'react';

interface DomainPermissionsComponentProps {
    className?: string;
}

export const DomainPermissionsComponent: React.FC<DomainPermissionsComponentProps> = ({
    className,
}) => {
    const [domainRoles] = useState<PermittedRoles[]>([]);
    const [domainUsers] = useState<PermittedUsers[]>([]);
    const [loading, setLoading] = useState(false);
    const [, setSelectedDomainId] = useState<number | null>(null);
    const { toast } = useToast();

    useEffect(() => {
        // Load initial data when component mounts
        loadDomainPermissions(1); // Default to domain ID 1
    });

    const loadDomainPermissions = async (domainId: number) => {
        try {
            setLoading(true);
            setSelectedDomainId(domainId);

            // Load both roles and users permissions for the domain
            // const [rolesResponse, usersResponse] = await Promise.all([
            //     RoleApiService.getDomainPermittedRoles(domainId),
            //     RoleApiService.getDomainPermittedUsers(domainId),
            // ]);

            // setDomainRoles(rolesResponse.roles || []);
            // setDomainUsers(usersResponse.users || []);
        } catch (error) {
            console.error('Error loading domain permissions:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить разрешения для домена',
                variant: 'destructive',
            });
        } finally {
            setLoading(false);
        }
    };

    if (loading) {
        return (
            <div className={className}>
                <Card>
                    <CardHeader>
                        <CardTitle>Разрешения домена</CardTitle>
                        <CardDescription>
                            Управление доступом ролей и пользователей к доменам
                        </CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className='space-y-4'>
                            <Skeleton className='h-4 w-full' />
                            <Skeleton className='h-4 w-3/4' />
                            <Skeleton className='h-4 w-1/2' />
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className={className}>
            <div className='space-y-6'>
                {/* Domain Roles Section */}
                <Card>
                    <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-4'>
                        <div>
                            <CardTitle className='flex items-center gap-2'>
                                <Users className='h-5 w-5' />
                                Роли с доступом к домену
                            </CardTitle>
                            <CardDescription>
                                Роли, которые имеют доступ к выбранному домену
                            </CardDescription>
                        </div>
                        <Button
                            onClick={() => {
                                // TODO: Implement add role dialog
                                toast({
                                    title: 'В разработке',
                                    description: 'Функция добавления ролей в разработке',
                                });
                            }}
                            size='sm'
                        >
                            <Plus className='h-4 w-4 mr-2' />
                            Добавить роль
                        </Button>
                    </CardHeader>
                    <CardContent>
                        {domainRoles.length === 0 ? (
                            <div className='text-center py-8 text-muted-foreground'>
                                <Users className='h-12 w-12 mx-auto mb-4 opacity-50' />
                                <p>Нет ролей с доступом к этому домену</p>
                            </div>
                        ) : (
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>ID</TableHead>
                                        <TableHead>Название роли</TableHead>
                                        <TableHead>Тип доступа</TableHead>
                                        <TableHead className='text-right'>Действия</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {domainRoles.map(() => (
                                        <TableRow>
                                            <TableCell className='font-mono'>
                                                {/* {role.role_id} */}
                                            </TableCell>
                                            <TableCell className='font-medium'>
                                                {/* {role.role_name} */}
                                            </TableCell>
                                            <TableCell>
                                                <Badge variant='outline'>
                                                    {/* {role.permission_type || 'read'} */}
                                                </Badge>
                                            </TableCell>
                                            <TableCell className='text-right'>
                                                <DropdownMenu>
                                                    <DropdownMenuTrigger asChild>
                                                        <Button
                                                            variant='ghost'
                                                            className='h-8 w-8 p-0'
                                                        >
                                                            <span className='sr-only'>
                                                                Открыть меню
                                                            </span>
                                                            <Edit className='h-4 w-4' />
                                                        </Button>
                                                    </DropdownMenuTrigger>
                                                    <DropdownMenuContent align='end'>
                                                        <DropdownMenuItem>
                                                            <Edit className='mr-2 h-4 w-4' />
                                                            Изменить доступ
                                                        </DropdownMenuItem>
                                                        <DropdownMenuItem className='text-destructive'>
                                                            <Trash2 className='mr-2 h-4 w-4' />
                                                            Удалить доступ
                                                        </DropdownMenuItem>
                                                    </DropdownMenuContent>
                                                </DropdownMenu>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        )}
                    </CardContent>
                </Card>

                {/* Domain Users Section */}
                <Card>
                    <CardHeader className='flex flex-row items-center justify-between space-y-0 pb-4'>
                        <div>
                            <CardTitle className='flex items-center gap-2'>
                                <User className='h-5 w-5' />
                                Пользователи с доступом к домену
                            </CardTitle>
                            <CardDescription>
                                Пользователи, которые имеют прямой доступ к выбранному домену
                            </CardDescription>
                        </div>
                        <Button
                            onClick={() => {
                                // TODO: Implement add user dialog
                                toast({
                                    title: 'В разработке',
                                    description: 'Функция добавления пользователей в разработке',
                                });
                            }}
                            size='sm'
                        >
                            <Plus className='h-4 w-4 mr-2' />
                            Добавить пользователя
                        </Button>
                    </CardHeader>
                    <CardContent>
                        {domainUsers.length === 0 ? (
                            <div className='text-center py-8 text-muted-foreground'>
                                <User className='h-12 w-12 mx-auto mb-4 opacity-50' />
                                <p>Нет пользователей с прямым доступом к этому домену</p>
                            </div>
                        ) : (
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>ID</TableHead>
                                        <TableHead>Имя пользователя</TableHead>
                                        <TableHead>Email</TableHead>
                                        <TableHead>Тип доступа</TableHead>
                                        <TableHead className='text-right'>Действия</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {domainUsers.map(() => (
                                        <TableRow>
                                            <TableCell className='font-mono'>
                                                {/* {user.user_id} */}
                                            </TableCell>
                                            <TableCell className='font-medium'>
                                                {/* {user.username} */}
                                            </TableCell>
                                            <TableCell>{/* {user.email} */}</TableCell>
                                            <TableCell>
                                                <Badge variant='outline'>
                                                    {/* {user.permission_type || 'read'} */}
                                                </Badge>
                                            </TableCell>
                                            <TableCell className='text-right'>
                                                <DropdownMenu>
                                                    <DropdownMenuTrigger asChild>
                                                        <Button
                                                            variant='ghost'
                                                            className='h-8 w-8 p-0'
                                                        >
                                                            <span className='sr-only'>
                                                                Открыть меню
                                                            </span>
                                                            <Edit className='h-4 w-4' />
                                                        </Button>
                                                    </DropdownMenuTrigger>
                                                    <DropdownMenuContent align='end'>
                                                        <DropdownMenuItem>
                                                            <Edit className='mr-2 h-4 w-4' />
                                                            Изменить доступ
                                                        </DropdownMenuItem>
                                                        <DropdownMenuItem className='text-destructive'>
                                                            <Trash2 className='mr-2 h-4 w-4' />
                                                            Удалить доступ
                                                        </DropdownMenuItem>
                                                    </DropdownMenuContent>
                                                </DropdownMenu>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        )}
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};

export default DomainPermissionsComponent;
