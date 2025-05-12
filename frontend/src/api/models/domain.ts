export interface TimeStamp {
    nanos: number;
    seconds: number;
}

export interface UpdateParams {
    cron: {
        dayOfMonth: number;
        dayOfWeek: number;
        hour: number;
        minute: number;
        month: number;
    };
    everyPeriod: number;
}

export interface Source {
    content: number[];
    createdAt: TimeStamp;
    credentials: number[];
    id: number;
    status: number;
    title: string;
    typ: number;
    updateParams: UpdateParams;
    updatedAt: TimeStamp;
    userId: number;
}

export interface SourcesResponse {
    sources: Source[];
}

export interface CreateSourceRequest {
    content: string;
    credentials?: string;
    title: string;
    typ: number;
    updateParams?: {
        cron?: {
            dayOfMonth: number;
            dayOfWeek: number;
            hour: number;
            minute: number;
            month: number;
        };
        everyPeriod?: number;
    };
}

export interface CreateDomainRequest {
    sourceIds: number[];
    title: string;
}

export interface Domain {
    createdAt: TimeStamp;
    id: number;
    sourceIds: number[];
    title: string;
    updatedAt: TimeStamp;
}

export interface DomainsResponse {
    domains: Domain[];
}

export interface ScenarioModelConfig {
    modelName: string;
    systemPrompt: string;
    temperature: number;
    topK: number;
    topP: number;
}

export interface ScenarioMultiQueryConfig {
    nQueries: number;
    queryModelName: string;
    useMultiquery: boolean;
}

export interface ScenarioRerankerConfig {
    rerankerMaxLength: number;
    rerankerModel: string;
    topK: number;
    useRerank: boolean;
}

export interface ScenarioVectorSearchConfig {
    searchByQuery: boolean;
    threshold: number;
    topN: number;
}

export interface CreateScenarioRequest {
    domainId: number;
    model: ScenarioModelConfig;
    multiQuery: ScenarioMultiQueryConfig;
    reranker: ScenarioRerankerConfig;
    vectorSearch: ScenarioVectorSearchConfig;
}

export enum SourceType {
    TypeWeb = 1,
    TypeSingleFile = 2,
    TypeArchivedFiles = 3,
    TypeWithCredentials = 4,
}
