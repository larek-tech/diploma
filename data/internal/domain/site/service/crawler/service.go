package crawler

type Service struct {
	httpClient   httpClient
	siteStore    siteStore
	pageStore    pageStore
	pageJobStore pageJobStore
	trManager    transactionalManager
}

func New(httpClient httpClient, siteStorage siteStore, pageStorage pageStore, pageJobStore pageJobStore, trManager transactionalManager) *Service {
	return &Service{
		httpClient:   httpClient,
		siteStore:    siteStorage,
		pageStore:    pageStorage,
		pageJobStore: pageJobStore,
		trManager:    trManager,
	}
}
