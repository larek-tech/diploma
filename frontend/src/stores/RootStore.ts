import ChatApiService from '@/api/ChatApiService';
import {DomainApiService} from '@/api/DomainApiService';
import {
  ChatSession,
  DeleteSessionParams,
  DisplayedChat,
  GetSessionParams,
  RenameSessionParams,
  SessionContentMessages,
  ShortChatSession,
  WSMessage,
  WSMessageType,
} from '@/api/models';
import {Domain, Scenario} from '@/api/models/domain';
import {LOCAL_STORAGE_KEY} from '@/auth/AuthProvider';
import {WS_URL} from '@/config';
import {makeAutoObservable, runInAction} from 'mobx';

export class RootStore {
    sessions: ShortChatSession[] = [];
    sessionsLoading: boolean = false;

    domains: Domain[] = [];
    domainsLoading: boolean = false;
    domainsOffset: number = 0;
    domainsLimit: number = 10;
    hasMoreDomains: boolean = true;
    selectedDomainId: number | null = null;
    selectedDomain: Domain | null = null;
    selectedScenarioId: number | null = null;
    selectedScenario: Scenario | null = null;

    activeSessionId: string | null = null;
    activeSession: ChatSession | null = null;
    activeDisplayedSession: DisplayedChat | null = null;
    activeSessionLoading: boolean = false;
    isChatDisabled: boolean = false;
    isModelAnswering: boolean = false;
    chatError: string | null = null;
    closedWebSocket: WebSocket | null = null;

    websocket: WebSocket | null = null;

    constructor() {
        makeAutoObservable(this);

        this.sessions = [];
        this.domains = [];
    }

    async setSelectedDomain(domainId: number) {
        this.selectedDomainId = domainId;

        // Получаем выбранный домен сетевым запросом
        const selectedDomain = await DomainApiService.getDomainById(domainId);

        this.selectedDomain = selectedDomain;

        const selectedScenarioId = selectedDomain.scenarioIds[0] || null;
        if (selectedScenarioId) {
            this.setSelectedScenario(selectedScenarioId);
        }
    }

    async setSelectedScenario(scenarioId: number) {
        this.selectedScenarioId = scenarioId;
        console.log('scenarioId', scenarioId);

        // Получаем выбранный сценарий сетевым запросом
        const selectedScenario = await DomainApiService.getScenarioById(scenarioId);

        this.selectedScenario = selectedScenario;
    }

    getSelectedDomainTitle(): string {
        if (!this.selectedDomain) return 'Выберите домен';

        return this.selectedDomain ? this.selectedDomain.title : 'Домен не найден';
    }

    getSelectedScenarioId(): number | null {
        return this.selectedScenarioId;
    }

    hasScenarios(): boolean {
        if (!this.selectedDomainId) return false;

        const selectedDomain = this.domains.find((domain) => domain.id === this.selectedDomainId);
        return selectedDomain ? selectedDomain.scenarioIds.length > 0 : false;
    }

    getAvailableScenarios(): number[] {
        if (!this.selectedDomainId) return [];

        const selectedDomain = this.domains.find((domain) => domain.id === this.selectedDomainId);
        return selectedDomain ? selectedDomain.scenarioIds : [];
    }

    async getSessions() {
        this.sessionsLoading = true;

        return ChatApiService.getSessions()
            .then(({ chats: sessions }) => {
                this.sessions = sessions;
            })
            .finally(() => {
                this.sessionsLoading = false;
            });
    }

    async getDomains(reset: boolean = false) {
        this.domainsLoading = true;

        if (reset) {
            this.domainsOffset = 0;
            this.domains = [];
            this.hasMoreDomains = true;
        }

        try {
            const response = await DomainApiService.getDomains(
                this.domainsOffset,
                this.domainsLimit
            );

            runInAction(() => {
                if (response.domains.length < this.domainsLimit) {
                    this.hasMoreDomains = false;
                }

                this.domains = [...this.domains, ...response.domains];
                this.domainsOffset += this.domainsLimit;
            });
        } catch (error) {
            console.error('Error loading domains:', error);
        } finally {
            runInAction(() => {
                this.domainsLoading = false;
            });
        }
    }

    async deleteSession({ id }: DeleteSessionParams) {
        return ChatApiService.deleteSession({ id }).then(() => {
            if (this.activeSessionId === id) {
                this.setActiveSessionId(null);
                this.activeSession = null;
            }
        });
    }

    async getSession({ id }: GetSessionParams) {
        this.activeSessionLoading = true;

        return ChatApiService.getSession({ id })
            .then((session) => {
                this.setActiveSession(session);
            })
            .finally(() => {
                this.activeSessionLoading = false;
            });
    }

    setActiveSession(session: ChatSession) {
        this.activeSession = session;

        this.activeDisplayedSession = {
            messages: session.content.map((message: SessionContentMessages) => ({
                query: message.query.content,
                response: message.response.content,
            })),
        };

        this.connectWebSocket(session.id);
    }

    setActiveSessionId(id: string | null) {
        if (id !== this.activeSessionId) {
            this.activeSessionId = id;
        }
    }

    renameSession({ id, title }: RenameSessionParams) {
        return ChatApiService.renameSession({ id, title });
    }

    async createSession() {
        return ChatApiService.createSession().then(async ({ id }) => {
            this.activeDisplayedSession = null;

            this.getSessions();

            this.connectWebSocket(id);
        });
    }

    private getAuthWSMessage(): WSMessage {
        const token = JSON.parse(localStorage.getItem(LOCAL_STORAGE_KEY) as string)?.user?.token;

        return {
            type: WSMessageType.Auth,
            content: token,
            isChunked: false,
            isLast: true,
            domainID: this.selectedDomainId || undefined,
            scenarioID: this.selectedScenarioId || undefined,
        };
    }

    connectWebSocket(sessionId: string) {
        this.disconnectWebSocket();

        this.websocket = new WebSocket(`${WS_URL}/${sessionId}`);

        this.websocket.onopen = () => {
            console.log('WebSocket connection opened');

            if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
                const wsMessage = this.getAuthWSMessage();

                this.websocket.send(JSON.stringify(wsMessage));

                this.setChatDisabled(false);
            }

            this.setActiveSessionId(sessionId);
        };

        this.websocket.onmessage = (event) => {
            const wsMessage: WSMessage = JSON.parse(event.data);

            runInAction(() => {
                const data = wsMessage.content;

                console.log(wsMessage);

                console.log(data);

                if (wsMessage.error) {
                    this.chatError = wsMessage.error;

                    // if (wsMessage.err === UNAUTHORIZED_ERR) {
                    //     localStorage.removeItem(LOCAL_STORAGE_KEY);

                    //     window.location.href = '/login';
                    // }

                    this.isModelAnswering = false;
                    this.isChatDisabled = false;
                }

                if (wsMessage.isChunked) {
                    this.isModelAnswering = true;
                    this.isChatDisabled = true;

                    this.processIncomingChunk(data || '');
                }

                // if (wsMessage.isChunked && !wsMessage.isLast) {
                //     // this.isModelAnswering = true;
                //     // this.isChatDisabled = true;

                //     this.processIncomingChunk(data as WSIncomingChunk);
                // } else if (!wsMessage.chunk && wsMessage.data && !wsMessage.data_type) {
                //     //!wsMessage.data_type значит, что это ответ модели (prediction или stock)
                //     this.processIncomingQuery(data as WSIncomingQuery);
                // }

                if (wsMessage.isLast || !wsMessage.isChunked) {
                    this.isModelAnswering = false;
                }

                if (wsMessage.isLast) {
                    this.isChatDisabled = false;
                }
            });
        };

        this.websocket.onclose = () => {
            console.log('WebSocket connection closed');

            this.isChatDisabled = true;
            this.closedWebSocket = this.websocket;

            this.reconnectWebSocket();
        };

        this.websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    sendMessage(message: WSMessage) {
        console.log('sendMessage', message);

        this.setIsModelAnswering(true);
        this.setChatDisabled(true);

        // Добавляем domainID и scenarioID в метаданные запроса, если они выбраны
        // message.queryMetadata = {
        //     ...message.queryMetadata,
        //     domainID: this.selectedDomainId || undefined,
        //     scenarioID: this.selectedScenarioId || undefined,
        // };

        message.domainID = this.selectedDomainId || undefined;
        message.scenarioID = this.selectedScenarioId || undefined;

        if (this.websocket && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(JSON.stringify(message));
        }

        if (this.isFirstMessageInSession()) {
            this.renameSession({
                id: this.activeSessionId as string,
                title: message.content?.slice(0, 60) || 'Без названия',
            });
        }

        this.addMessageToActiveSession(message);
    }

    disconnectWebSocket() {
        if (this.websocket) {
            this.setActiveSessionId(null);
            this.websocket.close();
        }
    }

    addMessageToActiveSession(message: WSMessage) {
        if (!this.activeSessionId) {
            return;
        }

        runInAction(() => {
            if (!this.activeDisplayedSession) {
                this.activeDisplayedSession = { messages: [] };
            }

            this.activeDisplayedSession?.messages.push({
                query: message.content || '',
                response: null,
            });
        });
    }

    // private processIncomingQuery(query: WSMessage) {
    //     console.log('processIncomingQuery', query);

    //     if (this.activeSessionId && this.activeDisplayedSession?.messages.length) {
    //         this.activeDisplayedSession.messages[
    //             this.activeDisplayedSession.messages.length - 1
    //         ].incomingMessage = {
    //             body: query.prompt,
    //             type: query.type as IncomingMessageType,
    //             status: query.status as IncomingMessageStatus,
    //             product: query.product,
    //             period: query.period,
    //         };
    //     }
    // }

    private processIncomingChunk(message: string) {
        if (this.activeSessionId && this.activeDisplayedSession?.messages.length) {
            const lastMessageIndex = this.activeDisplayedSession.messages.length - 1;
            const lastMessageBody = this.activeDisplayedSession.messages[lastMessageIndex].response;

            this.activeDisplayedSession.messages[lastMessageIndex].response = lastMessageBody
                ? lastMessageBody + message
                : message;
        }
    }

    setChatDisabled(isDisabled: boolean) {
        this.isChatDisabled = isDisabled;
    }

    setIsModelAnswering(isAnswering: boolean) {
        this.isModelAnswering = isAnswering;
    }

    cancelRequest() {
        // this.sendMessage({
        //     command: ChatCommand.Cancel,
        // });

        this.setChatDisabled(false);
        this.setIsModelAnswering(false);
    }

    private isFirstMessageInSession() {
        return !this.activeDisplayedSession?.messages.length;
    }

    private reconnectWebSocket() {
        if (this.activeSessionId) {
            this.connectWebSocket(this.activeSessionId);
        }
    }
}
