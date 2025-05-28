import {RoleForEntity} from '@/api/models/roles';
import {RoleApiService} from '@/api/RoleApiService';
import DomainManagementComponent from '@/components/DomainManagementComponent';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from '@/components/ui/card';
import {Tabs, TabsContent, TabsList, TabsTrigger} from '@/components/ui/tabs';
import {useToast} from '@/components/ui/use-toast';
import UserCreationComponent from '@/components/UserCreationComponent';
import {Globe, Plus, UserPlus} from 'lucide-react';
import {observer} from 'mobx-react-lite';
import {useEffect, useState} from 'react';

const RoleManagementPage = observer(() => {
    const [roles, setRoles] = useState<RoleForEntity[]>([]);
    const [, setLoading] = useState(false);
    const [, setIsCreateDialogOpen] = useState(false);
    const { toast } = useToast();

    useEffect(() => {
        loadRoles();
    }, []); // eslint-disable-line react-hooks/exhaustive-deps

    const loadRoles = async () => {
        try {
            setLoading(true);
            const response = await RoleApiService.listRoles(0, 100);
            setRoles(response.roles);
        } catch (error) {
            console.error('Error loading roles:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить список ролей',
                variant: 'destructive',
            });
        } finally {
            setLoading(false);
        }
    };

    // const handleCreateRole = async (data: CreateRoleRequest) => {
    //     try {
    //         const newRole = await RoleApiService.createRole(data);
    //         setRoles((prev) => [...prev, newRole]);
    //         setIsCreateDialogOpen(false);
    //         toast({
    //             title: 'Успех',
    //             description: `Роль "${data.name}" успешно создана`,
    //             variant: 'default',
    //         });
    //     } catch (error) {
    //         console.error('Error creating role:', error);
    //         toast({
    //             title: 'Ошибка',
    //             description: 'Не удалось создать роль',
    //             variant: 'destructive',
    //         });
    //     }
    // };

    // const handleDeleteRole = async (roleId: number) => {
    //     try {
    //         await RoleApiService.deleteRole(roleId);
    //         setRoles((prev) => prev.filter((role) => role.id !== roleId));
    //         toast({
    //             title: 'Успех',
    //             description: 'Роль успешно удалена',
    //             variant: 'default',
    //         });
    //     } catch (error) {
    //         console.error('Error deleting role:', error);
    //         toast({
    //             title: 'Ошибка',
    //             description: 'Не удалось удалить роль',
    //             variant: 'destructive',
    //         });
    //     }
    // };

    return (
        <div className='container mx-auto p-6 space-y-6'>
            <div className='flex items-center justify-between'>
                <div>
                    <h1 className='text-3xl font-bold tracking-tight'>Управление ролями</h1>
                    <p className='text-muted-foreground'>
                        Управляйте ролями пользователей и разрешениями для доменов и источников
                    </p>
                </div>
                <Button onClick={() => setIsCreateDialogOpen(true)}>
                    <Plus className='mr-2 h-4 w-4' />
                    Создать роль
                </Button>
            </div>

            <Tabs defaultValue='roles' className='space-y-4'>
                <TabsList>
                    <TabsTrigger value='user-management' className='flex items-center gap-2'>
                        <UserPlus className='h-4 w-4' />
                        Управление пользователями
                    </TabsTrigger>
                    <TabsTrigger value='domain-management' className='flex items-center gap-2'>
                        <Globe className='h-4 w-4' />
                        Управление доменами
                    </TabsTrigger>
                </TabsList>

                <TabsContent value='roles'>
                    <Card>
                        <CardHeader>
                            <CardTitle>Список ролей</CardTitle>
                            <CardDescription>
                                Управляйте ролями в системе. Создавайте, редактируйте и удаляйте
                                роли.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            {/* <RoleListComponent
                                roles={roles}
                                loading={loading}
                                onDelete={handleDeleteRole}
                                onRefresh={loadRoles}
                            /> */}
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value='user-roles'>
                    <Card>
                        <CardHeader>
                            <CardTitle>Назначение ролей пользователям</CardTitle>
                            <CardDescription>
                                Управляйте назначением ролей пользователям системы.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            {/* <UserRoleAssignmentComponent roles={roles} /> */}
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value='user-management'>
                    <Card>
                        <CardHeader>
                            <CardTitle>Управление пользователями</CardTitle>
                            <CardDescription>
                                Создавайте новых пользователей и просматривайте существующих.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <UserCreationComponent roles={roles} />
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value='domain-management'>
                    <Card>
                        <CardHeader>
                            <CardTitle>Управление доменами</CardTitle>
                            <CardDescription>
                                Управляйте доступом пользователей и ролей к доменам.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <DomainManagementComponent roles={roles} />
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value='domain-permissions'>
                    <Card>
                        <CardHeader>
                            <CardTitle>Разрешения доменов</CardTitle>
                            <CardDescription>
                                Управляйте доступом ролей и пользователей к доменам.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            {/* <DomainPermissionsComponent roles={roles} /> */}
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value='source-permissions'>
                    <Card>
                        <CardHeader>
                            <CardTitle>Разрешения источников</CardTitle>
                            <CardDescription>
                                Управляйте доступом ролей и пользователей к источникам данных.
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            {/* <SourcePermissionsComponent roles={roles} /> */}
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>

            {/* <CreateRoleDialog
                open={isCreateDialogOpen}
                onOpenChange={setIsCreateDialogOpen}
                onSubmit={handleCreateRole}
            /> */}
        </div>
    );
});

export default RoleManagementPage;
