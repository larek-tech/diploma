import { DomainForm } from '@/components/domain/DomainForm';

const CreateDomain = () => {
    return (
        <div className='container mx-auto py-10'>
            <h1 className='text-2xl font-bold mb-6'>Создание нового домена знаний</h1>
            <DomainForm />
        </div>
    );
};

export default CreateDomain;
