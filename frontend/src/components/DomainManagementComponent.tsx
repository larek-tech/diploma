import {Domain} from '@/api/models';
import {RoleForEntity, User} from '@/api/models/roles';
import {RoleApiService} from '@/api/RoleApiService';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {Checkbox} from '@/components/ui/checkbox';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import {useToast} from '@/components/ui/use-toast';
import {Globe, Search, Users} from 'lucide-react';
import {useEffect, useState} from 'react';

interface DomainManagementComponentProps {
    roles: RoleForEntity[];
}

export default function DomainManagementComponent({ roles }: DomainManagementComponentProps) {
    const [domains, setDomains] = useState<Domain[]>([]);
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedDomain, setSelectedDomain] = useState<Domain | null>(null);
    const [domainPermissions, setDomainPermissions] = useState<{
        userIds: number[];
        roleIds: number[];
    }>({ userIds: [], roleIds: [] });
    const { toast } = useToast();

    useEffect(() => {
        const loadData = async () => {
            try {
                setLoading(true);
                const [domainsResponse, usersResponse] = await Promise.all([
                    RoleApiService.getDomains(0, 100),
                    RoleApiService.listUsers(0, 100),
                ]);
                setDomains(domainsResponse.domains);
                setUsers(usersResponse.users);
            } catch (error) {
                console.error('Failed to load data:', error);
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось загрузить данные',
                    variant: 'destructive',
                });
            } finally {
                setLoading(false);
            }
        };

        loadData();
    }, [toast]);

    const loadDomainPermissions = async (domain: Domain) => {
        try {
            const [userPermissions, rolePermissions] = await Promise.all([
                RoleApiService.getDomainPermittedUsers(domain.id),
                RoleApiService.getDomainPermittedRoles(domain.id),
            ]);
            setDomainPermissions({
                userIds: userPermissions.userIds,
                roleIds: rolePermissions.roleIds,
            });
        } catch (error) {
            console.error('Error loading domain permissions:', error);
            setDomainPermissions({ userIds: [], roleIds: [] });
        }
    };

    const handleDomainSelect = async (domain: Domain) => {
        setSelectedDomain(domain);
        await loadDomainPermissions(domain);
    };

    const handleUserPermissionChange = (userId: number, checked: boolean) => {
        setDomainPermissions((prev) => ({
            ...prev,
            userIds: checked
                ? [...prev.userIds, userId]
                : prev.userIds.filter((id) => id !== userId),
        }));
    };

    const handleRolePermissionChange = (roleId: number, checked: boolean) => {
        setDomainPermissions((prev) => ({
            ...prev,
            roleIds: checked
                ? [...prev.roleIds, roleId]
                : prev.roleIds.filter((id) => id !== roleId),
        }));
    };

    const savePermissions = async () => {
        if (!selectedDomain) return;

        try {
            setLoading(true);
            await Promise.all([
                RoleApiService.updateDomainPermittedUsers(selectedDomain.id, {
                    resourceId: selectedDomain.id,
                    userIds: domainPermissions.userIds,
                }),
                RoleApiService.updateDomainPermittedRoles(selectedDomain.id, {
                    resourceId: selectedDomain.id,
                    roleIds: domainPermissions.roleIds,
                }),
            ]);

            toast({
                title: 'Успех',
                description: `Разрешения для домена "${selectedDomain.title}" успешно обновлены`,
                variant: 'default',
            });
        } catch (error) {
            console.error('Error saving permissions:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось сохранить разрешения',
                variant: 'destructive',
            });
        } finally {
            setLoading(false);
        }
    };

    const filteredDomains = domains.filter((domain) =>
        (domain.title?.toLowerCase() || '').includes((searchTerm || '').toLowerCase())
    );

    if (loading && domains.length === 0) {
        return <div className='p-4 text-center'>Загрузка...</div>;
    }

    return (
        <div className='space-y-6'>
            <div className='flex gap-6'>
                {/* Список доменов */}
                <Card className='w-1/2'>
                    <CardHeader>
                        <CardTitle className='flex items-center gap-2'>
                            <Globe className='h-5 w-5' />
                            Домены
                        </CardTitle>
                        <div className='relative'>
                            <Search className='absolute left-3 top-3 h-4 w-4 text-muted-foreground' />
                            <Input
                                placeholder='Поиск доменов...'
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className='pl-10'
                            />
                        </div>
                    </CardHeader>
                    <CardContent>
                        <div className='h-[400px] overflow-auto'>
                            <div className='space-y-2'>
                                {filteredDomains.map((domain) => (
                                    <div
                                        key={domain.id}
                                        className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                                            selectedDomain?.id === domain.id
                                                ? 'border-primary bg-primary/5'
                                                : 'border-border hover:bg-muted/50'
                                        }`}
                                        onClick={() => handleDomainSelect(domain)}
                                    >
                                        <div className='font-medium'>
                                            {domain.title || 'Untitled Domain'}
                                        </div>
                                        <div className='text-sm text-muted-foreground'>
                                            ID: {domain.id || 'N/A'} • Источников:{' '}
                                            {domain.sourceIds?.length || 0}
                                        </div>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </CardContent>
                </Card>

                {/* Управление разрешениями */}
                <Card className='w-1/2'>
                    <CardHeader>
                        <CardTitle className='flex items-center gap-2'>
                            <Users className='h-5 w-5' />
                            Управление доступом
                            {selectedDomain && (
                                <Badge variant='outline'>{selectedDomain.title}</Badge>
                            )}
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        {selectedDomain ? (
                            <div className='space-y-6'>
                                {/* Пользователи */}
                                <div>
                                    <Label className='text-base font-medium'>
                                        Разрешенные пользователи
                                    </Label>
                                    <div className='h-[150px] mt-2 overflow-auto'>
                                        <div className='space-y-2'>
                                            {users.map((user) => (
                                                <div
                                                    key={user.id}
                                                    className='flex items-center space-x-2'
                                                >
                                                    <Checkbox
                                                        id={`user-${user.id}`}
                                                        checked={domainPermissions.userIds.includes(
                                                            user.id
                                                        )}
                                                        onCheckedChange={(checked) =>
                                                            handleUserPermissionChange(
                                                                user.id,
                                                                !!checked
                                                            )
                                                        }
                                                    />
                                                    <Label
                                                        htmlFor={`user-${user.id}`}
                                                        className='text-sm cursor-pointer'
                                                    >
                                                        {user.firstName || ''} {user.lastName || ''}{' '}
                                                        ({user.login || 'N/A'})
                                                    </Label>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </div>

                                {/* Роли */}
                                <div>
                                    <Label className='text-base font-medium'>
                                        Разрешенные роли
                                    </Label>
                                    <div className='h-[150px] mt-2 overflow-auto'>
                                        <div className='space-y-2'>
                                            {roles.map((role) => (
                                                <div
                                                    key={role.id}
                                                    className='flex items-center space-x-2'
                                                >
                                                    <Checkbox
                                                        id={`role-${role.id}`}
                                                        checked={domainPermissions.roleIds?.includes(
                                                            role.id
                                                        )}
                                                        onCheckedChange={(checked) =>
                                                            handleRolePermissionChange(
                                                                role.id,
                                                                !!checked
                                                            )
                                                        }
                                                    />
                                                    <Label
                                                        htmlFor={`role-${role.id}`}
                                                        className='text-sm cursor-pointer'
                                                    >
                                                        {role.name}
                                                    </Label>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                </div>

                                <Button
                                    onClick={savePermissions}
                                    disabled={loading}
                                    className='w-full'
                                >
                                    {loading ? 'Сохранение...' : 'Сохранить разрешения'}
                                </Button>
                            </div>
                        ) : (
                            <div className='text-center text-muted-foreground py-8'>
                                Выберите домен для управления разрешениями
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
