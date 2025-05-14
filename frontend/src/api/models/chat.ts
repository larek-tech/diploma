import { Scenario } from './domain';
import { Time } from './time';

export interface ShortChatSession {
    createdAt: Time;
    id: string;
    title: string;
    updatedAt: Time;
    userId: number;
}

export interface ChatSession {
    createdAt: Time;
    id: string;
    content: SessionContentMessages[];
    title: string;
    updatedAt: Time;
    userId: number;
}

export interface SessionContentMessages {
    query: ChatSessionQuery;
    response: ChatSessionResponse;
}

export interface ChatSessionQuery {
    chatId: string;
    content: string;
    createdAt: Time;
    domainId: number;
    id: number;
    metadata: number[];
    scenarioId: number;
    sourceIds: number[];
    userId: number;
}

export interface ChatSessionResponse {
    chatId: string;
    content: string;
    createdAt: Time;
    id: number;
    metadata: number[];
    queryId: number;
    status: number;
    updatedAt: Time;
}

export interface DisplayedChat {
    messages: ChatConversation[];
}

export interface ChatConversation {
    query: string;
    response: string | null;
}

export interface GetSessionsResponse {
    chats: ShortChatSession[];
}

export interface CreateSessionResponse {
    id: string;
}

export interface GetSessionParams {
    id: string;
}

export interface RenameSessionParams {
    id: string;
    title: string;
}

export interface DeleteSessionParams {
    id: string;
}

export interface WSMessage {
    type: WSMessageType;
    content?: string;
    isChunked?: boolean;
    isLast?: boolean;
    sourceIDs?: string[];
    domainID?: number;
    scenarioID?: number;
    queryMetadata?: QueryMetadata;
    error?: string;
}

export enum WSMessageType {
    Auth = 'auth',
    Query = 'query',
    Chunk = 'chunk',
    Error = 'error',
}

export interface QueryMetadata {
    domainID?: number;
    scenarioID?: number;
    scenario?: Scenario;
}

// // SocketMessageType enum defining type of socket message.
// type SocketMessageType string

// const (
//  // TypeAuth content is auth token (without Bearer).
//  TypeAuth  SocketMessageType = "auth"
//  // TypeQuery content is general message.
//  TypeQuery SocketMessageType = "query"
//  // TypeChunk content is chunked response.
//  TypeChunk SocketMessageType = "chunk"
//  // TypeError content is empty, got error, chat is stopped.
//  TypeError SocketMessageType = "error"
// )

// // SocketMessage is a model for incoming and outgoing messages for websocket.
// type SocketMessage struct {
//  Type          SocketMessageType `json:"type"`
//  Content       string            `json:"content,omitempty"`
//  IsChunked     bool              `json:"isChunked"`
//  IsLast        bool              `json:"isLast"`
//  SourceIDs     []string          `json:"sourceIDs,omitempty"`
//  QueryMetadata QueryMetadata     `json:"queryMetadata,omitempty"`
//  Err           string            `json:"error,omitempty"`
// }

// // QueryMetadata stores information about chosen domain, sources and scenario for query.
// type QueryMetadata struct {
//  DomainID *int64       `json:"domainID,omitempty"`
//  Scenario *pb.Scenario `json:"scenario,omitempty"`
// }

// export interface ChatQuery {
//     id: number;
//     prompt: string;
//     product: string;
//     period?: string;
//     type: string;
//     status: string;
//     created_at: string;
// }

// export interface ChatResponse {
//     created_at: string;
//     body: string;
//     status: string;
//     data: PredictionResponse | StockResponse;
//     data_type: ModelResponseType;
// }

// export interface SessionContent {
//     query: ChatQuery;
//     response: ChatResponse;
// }

// export interface ChatSession {
//     id: string;
//     title: string;
//     content: SessionContent[];
//     editable: boolean;
//     tg: boolean;
// }

// export interface WSOutcomingMessage {
//     prompt?: string;
//     command?: ChatCommand;
//     period?: string;
//     product?: string;
// }

// export enum ChatCommand {
//     Valid = 'valid',
//     Invalid = 'invalid',
//     Cancel = 'cancel',
// }

// export enum IncomingMessageType {
//     Stock = 'STOCK',
//     Prediction = 'PREDICTION',
//     Undefined = 'UNDEFINED',
// }

// export enum IncomingMessageStatus {
//     Pending = 'PENDING',
//     Valid = 'VALID',
//     Invalid = 'INVALID',
// }

// export interface WSMessage {
//     data: WSIncomingQuery | WSIncomingChunk | PredictionResponse | StockResponse;
//     finish: boolean;
//     chunk: boolean;
//     err?: string;
//     data_type?: ModelResponseType;
// }

// export interface WSIncomingQuery {
//     created_at: string;
//     prompt: string;
//     period: string;
//     product: string;
//     type: IncomingMessageType;
//     status: string;
//     id: number;
// }

// export interface WSIncomingChunk {
//     info: string;
// }

// export interface ChatConversation {
//     outcomingMessage?: DisplayedOutcomingMessage;
//     incomingMessage?: DisplayedIncomingMessage;
// }

// export interface DisplayedOutcomingMessage {
//     prompt: string;
// }

// export interface DisplayedIncomingMessage {
//     type: IncomingMessageType;
//     status: IncomingMessageStatus;
//     body: string;
//     product?: string;
//     period?: string;
//     prediction?: { forecast: PurchasePlan[]; history: PurchasePlan[] };
//     stocks?: StockResponse['data'];
//     outputJson?: OutputJson;
// }

// export const UNAUTHORIZED_ERR = 'invalid JWT';
