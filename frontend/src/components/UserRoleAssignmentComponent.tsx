import {Source} from '@/api/models';
import {RoleForEntity, User} from '@/api/models/roles';
import {RoleApiService} from '@/api/RoleApiService';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card';
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from '@/components/ui/select';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from '@/components/ui/table';
import {useToast} from '@/components/ui/use-toast';
import {Loader2, RefreshCw, UserMinus} from 'lucide-react';
import React, {useEffect, useState} from 'react';

interface UserRoleAssignmentComponentProps {
    roles: RoleForEntity[];
}

const UserRoleAssignmentComponent: React.FC<UserRoleAssignmentComponentProps> = ({ roles }) => {
    const [users, setUsers] = useState<User[]>([]);
    const [sources, setSources] = useState<Source[]>([]);
    const [loading, setLoading] = useState(false);
    const [assigningRole, setAssigningRole] = useState<{ userId: number; roleId: number } | null>(
        null
    );
    const [removingRole, setRemovingRole] = useState<{ userId: number; roleId: number } | null>(
        null
    );
    const { toast } = useToast();

    useEffect(() => {
        loadData();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    const loadData = async () => {
        try {
            setLoading(true);
            const [usersResponse, sourcesResponse] = await Promise.all([
                RoleApiService.listUsers(0, 100),
                RoleApiService.getSources(0, 100),
            ]);
            setUsers(usersResponse.users);
            setSources(sourcesResponse.sources);
        } catch (error) {
            console.error('Error loading data:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить данные',
                variant: 'destructive',
            });
        } finally {
            setLoading(false);
        }
    };

    const handleAssignRole = async (userId: number, roleId: number) => {
        try {
            setAssigningRole({ userId, roleId });
            await RoleApiService.updateUserRole(userId, roleId);

            // Update user roles in local state
            setUsers((prev) =>
                prev.map((user) => {
                    if (user.id === userId) {
                        const userRoles = user.roles || [];
                        if (!userRoles.includes(roleId)) {
                            return {
                                ...user,
                                roles: [...userRoles, roleId],
                            };
                        }
                    }
                    return user;
                })
            );

            toast({
                title: 'Успех',
                description: 'Роль успешно назначена пользователю',
                variant: 'default',
            });
        } catch (error) {
            console.error('Error assigning role:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось назначить роль пользователю',
                variant: 'destructive',
            });
        } finally {
            setAssigningRole(null);
        }
    };

    const handleRemoveRole = async (userId: number, roleId: number) => {
        try {
            setRemovingRole({ userId, roleId });
            await RoleApiService.removeUserRole(userId, roleId);

            // Update user roles in local state
            setUsers((prev) =>
                prev.map((user) => {
                    if (user.id === userId) {
                        const userRoles = user.roles || [];
                        return {
                            ...user,
                            roles: userRoles.filter((id) => id !== roleId),
                        };
                    }
                    return user;
                })
            );

            toast({
                title: 'Успех',
                description: 'Роль успешно удалена у пользователя',
                variant: 'default',
            });
        } catch (error) {
            console.error('Error removing role:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось удалить роль у пользователя',
                variant: 'destructive',
            });
        } finally {
            setRemovingRole(null);
        }
    };

    const getAvailableRoles = (user: User) => {
        const userRoles = user.roles || [];
        return roles.filter((role) => !userRoles.includes(role.id));
    };

    const getUserRoleName = (roleId: number) => {
        const role = roles.find((r) => r.id === roleId);
        return role ? role.name : `Role ${roleId}`;
    };

    if (loading) {
        return (
            <div className='flex items-center justify-center p-8'>
                <Loader2 className='h-8 w-8 animate-spin' />
                <span className='ml-2'>Загрузка данных...</span>
            </div>
        );
    }

    return (
        <div className='space-y-6'>
            <div className='flex items-center justify-between'>
                <div>
                    <h3 className='text-lg font-medium'>Пользователи и их роли</h3>
                    <p className='text-sm text-muted-foreground'>
                        Управляйте назначением ролей для каждого пользователя
                    </p>
                </div>
                <Button onClick={loadData} variant='outline' size='sm'>
                    <RefreshCw className='h-4 w-4 mr-2' />
                    Обновить
                </Button>
            </div>

            <Card>
                <CardContent className='p-0'>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Пользователь</TableHead>
                                <TableHead>Email</TableHead>
                                <TableHead>Текущие роли</TableHead>
                                <TableHead>Назначить роль</TableHead>
                                <TableHead>Действия</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {users.map((user) => {
                                const availableRoles = getAvailableRoles(user);
                                return (
                                    <TableRow key={user.id}>
                                        <TableCell>
                                            <div>
                                                <div className='font-medium'>
                                                    {user.firstName} {user.lastName}
                                                </div>
                                                <div className='text-sm text-muted-foreground'>
                                                    @{user.login}
                                                </div>
                                            </div>
                                        </TableCell>
                                        <TableCell>{user.email}</TableCell>
                                        <TableCell>
                                            <div className='flex flex-wrap gap-1'>
                                                {(user.roles || []).map((roleId) => (
                                                    <Badge
                                                        key={roleId}
                                                        variant='secondary'
                                                        className='flex items-center gap-1'
                                                    >
                                                        {getUserRoleName(roleId)}
                                                        <Button
                                                            size='sm'
                                                            variant='ghost'
                                                            className='h-auto p-0 ml-1 text-destructive hover:text-destructive'
                                                            onClick={() =>
                                                                handleRemoveRole(user.id, roleId)
                                                            }
                                                            disabled={
                                                                removingRole?.userId === user.id &&
                                                                removingRole?.roleId === roleId
                                                            }
                                                        >
                                                            {removingRole?.userId === user.id &&
                                                            removingRole?.roleId === roleId ? (
                                                                <Loader2 className='h-3 w-3 animate-spin' />
                                                            ) : (
                                                                <UserMinus className='h-3 w-3' />
                                                            )}
                                                        </Button>
                                                    </Badge>
                                                ))}
                                                {(!user.roles || user.roles.length === 0) && (
                                                    <span className='text-sm text-muted-foreground'>
                                                        Нет назначенных ролей
                                                    </span>
                                                )}
                                            </div>
                                        </TableCell>
                                        <TableCell>
                                            {availableRoles.length > 0 ? (
                                                <Select
                                                    onValueChange={(value) =>
                                                        handleAssignRole(user.id, parseInt(value))
                                                    }
                                                    disabled={assigningRole?.userId === user.id}
                                                >
                                                    <SelectTrigger className='w-40'>
                                                        <SelectValue placeholder='Выберите роль' />
                                                    </SelectTrigger>
                                                    <SelectContent>
                                                        {availableRoles.map((role) => (
                                                            <SelectItem
                                                                key={role.id}
                                                                value={role.id.toString()}
                                                            >
                                                                {role.name}
                                                            </SelectItem>
                                                        ))}
                                                    </SelectContent>
                                                </Select>
                                            ) : (
                                                <span className='text-sm text-muted-foreground'>
                                                    Все роли назначены
                                                </span>
                                            )}
                                        </TableCell>
                                        <TableCell>
                                            {assigningRole?.userId === user.id && (
                                                <div className='flex items-center gap-2'>
                                                    <Loader2 className='h-4 w-4 animate-spin' />
                                                    <span className='text-sm'>Назначение...</span>
                                                </div>
                                            )}
                                        </TableCell>
                                    </TableRow>
                                );
                            })}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            {sources.length > 0 && (
                <Card>
                    <CardHeader>
                        <CardTitle>Доступные источники</CardTitle>
                        <CardDescription>Список источников данных в системе</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'>
                            {sources.map((source) => (
                                <Card key={source.id} className='p-4'>
                                    <div className='font-medium'>{source.title}</div>
                                    <div className='text-sm text-muted-foreground mt-1'>
                                        Тип: {source.typ}
                                    </div>
                                </Card>
                            ))}
                        </div>
                    </CardContent>
                </Card>
            )}
        </div>
    );
};

export default UserRoleAssignmentComponent;
