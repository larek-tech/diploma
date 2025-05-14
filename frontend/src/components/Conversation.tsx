import UserMessage from './UserMessage';
import ModelMessage from './ModelMessage';
import { observer } from 'mobx-react-lite';
import { ChatConversation } from '@/api/models';

type Props = {
    conversation: ChatConversation;
    isLastConversation: boolean;
};

const Conversation = observer(({ conversation, isLastConversation }: Props) => {
    return (
        <div className='w-full'>
            {conversation.query && <UserMessage message={conversation.query} />}

            {conversation.response && (
                <ModelMessage
                    incomingMessage={conversation.response}
                    isLastMessage={isLastConversation}
                />
            )}
        </div>
    );
});

export default Conversation;
