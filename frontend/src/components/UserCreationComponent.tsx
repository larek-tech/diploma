import {CreateUserRequest, RoleForEntity, User} from '@/api/models/roles';
import {RoleApiService} from '@/api/RoleApiService';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from '@/components/ui/table';
import {useToast} from '@/components/ui/use-toast';
import {Plus, Search, Users} from 'lucide-react';
import React, {useEffect, useState} from 'react';

interface UserCreationComponentProps {
    roles: RoleForEntity[];
}

export default function UserCreationComponent({ roles }: UserCreationComponentProps) {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);
    const [searchTerm, setSearchTerm] = useState('');
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [formData, setFormData] = useState<CreateUserRequest>({
        email: '',
        password: '',
    });
    const { toast } = useToast();

    useEffect(() => {
        const loadUsers = async () => {
            try {
                setLoading(true);
                const response = await RoleApiService.listUsers(0, 100);
                setUsers(response.users);
            } catch (error) {
                console.error('Failed to load users:', error);
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось загрузить пользователей',
                    variant: 'destructive',
                });
            } finally {
                setLoading(false);
            }
        };

        loadUsers();
    }, [toast]);

    const handleCreateUser = async (e: React.FormEvent) => {
        e.preventDefault();

        // Валидация
        if (!formData.email || !formData.password) {
            toast({
                title: 'Ошибка',
                description: 'Email и пароль обязательны для заполнения',
                variant: 'destructive',
            });
            return;
        }

        if (!formData.email.includes('@')) {
            toast({
                title: 'Ошибка',
                description: 'Введите корректный email адрес',
                variant: 'destructive',
            });
            return;
        }

        if (formData.password.length < 6) {
            toast({
                title: 'Ошибка',
                description: 'Пароль должен содержать минимум 6 символов',
                variant: 'destructive',
            });
            return;
        }

        try {
            setLoading(true);
            const newUser = await RoleApiService.createUser(formData);
            setUsers((prev) => [...prev, newUser]);
            setIsCreateDialogOpen(false);
            setFormData({
                email: '',
                password: '',
            });

            toast({
                title: 'Успех',
                description: `Пользователь "${formData.email}" успешно создан`,
                variant: 'default',
            });
        } catch (error) {
            console.error('Error creating user:', error);
            toast({
                title: 'Ошибка',
                description: 'Не удалось создать пользователя',
                variant: 'destructive',
            });
        } finally {
            setLoading(false);
        }
    };

    const handleInputChange = (field: keyof CreateUserRequest, value: string) => {
        setFormData((prev) => ({
            ...prev,
            [field]: value,
        }));
    };

    const formatDate = (timestamp: { seconds: number; nanos: number }) => {
        return new Date(timestamp.seconds * 1000).toLocaleDateString('ru-RU');
    };

    const filteredUsers = users.filter(
        (user) =>
            (user.firstName?.toLowerCase() || '').includes((searchTerm || '').toLowerCase()) ||
            (user.lastName?.toLowerCase() || '').includes((searchTerm || '').toLowerCase()) ||
            (user.login?.toLowerCase() || '').includes((searchTerm || '').toLowerCase()) ||
            (user.email?.toLowerCase() || '').includes((searchTerm || '').toLowerCase())
    );

    return (
        <div className='space-y-6'>
            <Card>
                <CardHeader>
                    <div className='flex items-center justify-between'>
                        <CardTitle className='flex items-center gap-2'>
                            <Users className='h-5 w-5' />
                            Управление пользователями
                        </CardTitle>
                        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                            <DialogTrigger asChild>
                                <Button>
                                    <Plus className='mr-2 h-4 w-4' />
                                    Создать пользователя
                                </Button>
                            </DialogTrigger>
                            <DialogContent className='sm:max-w-md'>
                                <form onSubmit={handleCreateUser}>
                                    <DialogHeader>
                                        <DialogTitle>Создание нового пользователя</DialogTitle>
                                        <DialogDescription>
                                            Заполните все поля для создания нового пользователя
                                        </DialogDescription>
                                    </DialogHeader>
                                    <div className='space-y-4 py-4'>
                                        <div className='space-y-2'>
                                            <Label htmlFor='email'>Email</Label>
                                            <Input
                                                id='email'
                                                type='email'
                                                value={formData.email}
                                                onChange={(e) =>
                                                    handleInputChange('email', e.target.value)
                                                }
                                                placeholder='Введите email'
                                                required
                                            />
                                        </div>
                                        <div className='space-y-2'>
                                            <Label htmlFor='password'>Пароль</Label>
                                            <Input
                                                id='password'
                                                type='password'
                                                value={formData.password}
                                                onChange={(e) =>
                                                    handleInputChange('password', e.target.value)
                                                }
                                                placeholder='Введите пароль (минимум 6 символов)'
                                                minLength={6}
                                                required
                                            />
                                        </div>
                                    </div>
                                    <DialogFooter>
                                        <Button
                                            type='button'
                                            variant='outline'
                                            onClick={() => setIsCreateDialogOpen(false)}
                                        >
                                            Отмена
                                        </Button>
                                        <Button type='submit' disabled={loading}>
                                            {loading ? 'Создание...' : 'Создать'}
                                        </Button>
                                    </DialogFooter>
                                </form>
                            </DialogContent>
                        </Dialog>
                    </div>
                    <div className='relative'>
                        <Search className='absolute left-3 top-3 h-4 w-4 text-muted-foreground' />
                        <Input
                            placeholder='Поиск пользователей...'
                            value={searchTerm}
                            onChange={(e) => setSearchTerm(e.target.value)}
                            className='pl-10'
                        />
                    </div>
                </CardHeader>
                <CardContent>
                    {loading && users.length === 0 ? (
                        <div className='text-center py-8'>Загрузка...</div>
                    ) : (
                        <div className='h-[500px] overflow-auto'>
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>ID</TableHead>
                                        <TableHead>Логин</TableHead>
                                        <TableHead>Имя</TableHead>
                                        <TableHead>Email</TableHead>
                                        <TableHead>Роли</TableHead>
                                        <TableHead>Дата создания</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {filteredUsers.map((user) => (
                                        <TableRow key={user.id}>
                                            <TableCell className='font-mono text-sm'>
                                                {user.id || 'N/A'}
                                            </TableCell>
                                            <TableCell className='font-medium'>
                                                {user.login || 'N/A'}
                                            </TableCell>
                                            <TableCell>
                                                {user.firstName || ''} {user.lastName || ''}
                                            </TableCell>
                                            <TableCell>{user.email || 'N/A'}</TableCell>
                                            <TableCell>
                                                {user.roles && user.roles.length > 0 ? (
                                                    <div className='flex flex-wrap gap-1'>
                                                        {user.roles.map((roleId) => {
                                                            const role = roles.find(
                                                                (r) => r.id === roleId
                                                            );
                                                            return (
                                                                <Badge
                                                                    key={roleId}
                                                                    variant='secondary'
                                                                    className='text-xs'
                                                                >
                                                                    {role?.name || `Роль ${roleId}`}
                                                                </Badge>
                                                            );
                                                        })}
                                                    </div>
                                                ) : (
                                                    <span className='text-muted-foreground'>
                                                        Нет ролей
                                                    </span>
                                                )}
                                            </TableCell>
                                            <TableCell className='text-sm text-muted-foreground'>
                                                {formatDate(user.createdAt)}
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
