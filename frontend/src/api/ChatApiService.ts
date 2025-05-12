import {
    ChatSession,
    CreateSessionResponse,
    DeleteSessionParams,
    GetSessionParams,
    GetSessionsResponse,
    RenameSessionParams,
} from './models';
import { get, post, del, put } from './http';

class ChatApiService {
    public async createSession() {
        const response = await post<CreateSessionResponse>('/api/v1/chat', {});

        return response;
    }

    public async getSessions() {
        const response = await get<GetSessionsResponse>('/api/v1/chat/list');

        return response;
    }

    public async renameSession({ id, title }: RenameSessionParams) {
        await put(`/chat/session/rename`, { id, title });
    }

    public async deleteSession({ id }: DeleteSessionParams) {
        await del(`/chat/session/${id}`);
    }

    public async getSession({ id }: GetSessionParams) {
        const response = await get<ChatSession>(`/chat/session/${id}`);

        return response;
    }
}

export default new ChatApiService();
