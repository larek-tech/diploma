import axiosInstance from './axiosInstance';
import {
  CreateDomainRequest,
  CreateSourceRequest,
  Domain,
  DomainsResponse,
  Scenario,
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

    static async getDomains(offset: number = 0, limit: number = 10): Promise<DomainsResponse> {
        const response = await axiosInstance.get<DomainsResponse>('/api/v1/domain/list', {
            params: { offset, limit },
        });
        return response.data;
    }

    static async getDomainById(domainId: number): Promise<Domain> {
        const response = await axiosInstance.get<Domain>(`/api/v1/domain/${domainId}`);

        return response.data;
    }

    static async createScenario(data: Scenario): Promise<Scenario> {
        const response = await axiosInstance.post<Scenario>('/api/v1/scenario', data);
        return response.data;
    }

    static async getScenarios(
        offset: number = 0,
        limit: number = 1000
    ): Promise<{ scenarios: Scenario[] }> {
        const response = await axiosInstance.get<{ scenarios: Scenario[] }>(
            '/api/v1/scenario/list',
            {
                params: { offset, limit },
            }
        );
        return response.data;
    }

    static async getScenarioById(scenarioId: number): Promise<Scenario> {
        const response = await axiosInstance.get<Scenario>(`/api/v1/scenario/${scenarioId}`);

        return response.data;
    }
}
