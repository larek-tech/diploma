import { Button } from '@/components/ui/button';
import { useNavigate } from 'react-router-dom';
import { Pages } from '@/router/constants';
import { PlusCircle } from 'lucide-react';

const HomeBanner = () => {
    const navigate = useNavigate();
    return (
        <div className='flex w-full justify-center items-center min-h-[60vh]'>
            <div className='space-y-6 text-center max-w-md'>
                <PlusCircle className='mx-auto h-12 w-12 text-primary mb-2' />
                <h2 className='text-2xl font-bold mb-2'>
                    Для продолжения выберете домен или создайте новый
                </h2>
                <p className='text-gray-500 mb-4'>
                    Создайте свой первый домен, чтобы начать работу с чатами и сценариями.
                </p>
                <Button size='lg' onClick={() => navigate(`/${Pages.CreateDomain}`)}>
                    Создать домен
                </Button>
            </div>
        </div>
    );
};

export default HomeBanner;
