import {useParams} from 'react-router-dom';

const ChatFromHistory = () => {
    const { sessionId } = useParams();

    return (
        <div>
            <h1>Chat From History</h1>
        </div>
    );
};

export default ChatFromHistory;
