import { get, post, del, put } from './http';
import {
    DeleteSessionParams,
    GetSessionParams,
    RenameSessionParams,
    ChatSession,
    GetSessionsResponse,
} from './models';

class ChatApiService {
    public async createSession() {
        const response = await post<ChatSession>('/api/v1/chat', {});

        return response;
    }

    public async getSessions() {
        const response = await get<GetSessionsResponse>('/api/v1/chat/list');

        return response;
    }

    public async renameSession({ id, title }: RenameSessionParams) {
        await put(`api/v1/chat/${id}`, { title });
    }

    public async deleteSession({ id }: DeleteSessionParams) {
        await del(`api/v1/chat/${id}`);
    }

    public async getSession({ id }: GetSessionParams) {
        const response = await get<ChatSession>(`/api/v1/chat/history/${id}`);

        return response;
    }
}

export default new ChatApiService();
