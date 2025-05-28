import {CreateRoleRequest} from '@/api/models/roles';
import {Button} from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import {Input} from '@/components/ui/input';
import {Label} from '@/components/ui/label';
import React, {useState} from 'react';

interface CreateRoleDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onSubmit: (data: CreateRoleRequest) => Promise<void>;
}

export const CreateRoleDialog: React.FC<CreateRoleDialogProps> = ({
    open,
    onOpenChange,
    onSubmit,
}) => {
    const [formData, setFormData] = useState<CreateRoleRequest>({
        name: '',
    });
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!formData.name.trim()) {
            return;
        }

        try {
            setIsSubmitting(true);
            await onSubmit(formData);
            // Reset form after successful submission
            setFormData({ name: '' });
        } catch (error) {
            console.error('Error creating role:', error);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleCancel = () => {
        setFormData({ name: '' });
        onOpenChange(false);
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className='sm:max-w-[425px]'>
                <form onSubmit={handleSubmit}>
                    <DialogHeader>
                        <DialogTitle>Создать новую роль</DialogTitle>
                        <DialogDescription>
                            Создайте новую роль для управления доступом пользователей к ресурсам.
                        </DialogDescription>
                    </DialogHeader>

                    <div className='grid gap-4 py-4'>
                        <div className='grid grid-cols-4 items-center gap-4'>
                            <Label htmlFor='name' className='text-right'>
                                Название*
                            </Label>
                            <Input
                                id='name'
                                value={formData.name}
                                onChange={(e) =>
                                    setFormData((prev) => ({ ...prev, name: e.target.value }))
                                }
                                className='col-span-3'
                                placeholder='Введите название роли'
                                required
                                disabled={isSubmitting}
                            />
                        </div>
                    </div>

                    <DialogFooter>
                        <Button
                            type='button'
                            variant='outline'
                            onClick={handleCancel}
                            disabled={isSubmitting}
                        >
                            Отмена
                        </Button>
                        <Button type='submit' disabled={!formData.name.trim() || isSubmitting}>
                            {isSubmitting ? 'Создание...' : 'Создать роль'}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    );
};

export default CreateRoleDialog;
