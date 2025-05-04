package crawler

type Service struct {
	httpClient httpClient
	siteStore  siteStore
	pageStore  pageStore
	trManager  transactionalManager
}

func New(httpClient httpClient, siteStorage siteStore, pageStorage pageStore, trManager transactionalManager) *Service {
	return &Service{
		httpClient: httpClient,
		siteStore:  siteStorage,
		pageStore:  pageStorage,
		trManager:  trManager,
	}
}
