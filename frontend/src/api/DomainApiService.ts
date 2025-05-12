import axiosInstance from './axiosInstance';
import {
    CreateDomainRequest,
    CreateScenarioRequest,
    CreateSourceRequest,
    Domain,
    Source,
    SourcesResponse,
} from './models';

export class DomainApiService {
    static async getSources(): Promise<SourcesResponse> {
        const response = await axiosInstance.get<SourcesResponse>('/api/v1/source/list');
        return response.data;
    }

    static async createSource(data: CreateSourceRequest): Promise<Source> {
        const response = await axiosInstance.post<Source>('/api/v1/source', data);
        return response.data;
    }

    static async createDomain(data: CreateDomainRequest): Promise<Domain> {
        const response = await axiosInstance.post<Domain>('/api/v1/domain', data);
        return response.data;
    }

    static async createScenario(data: CreateScenarioRequest): Promise<unknown> {
        const response = await axiosInstance.post<unknown>('/api/v1/scenario', data);
        return response.data;
    }
}
