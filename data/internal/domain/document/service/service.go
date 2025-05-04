package service

type Service struct {
	documentStorage documentStorage
	chunkStorage    chunkStorage
	embedder        embedder
	trManager       trManager
}

func New(documentStorage documentStorage, chunkStorage chunkStorage, embedder embedder, trManager trManager) *Service {
	return &Service{
		documentStorage: documentStorage,
		chunkStorage:    chunkStorage,
		embedder:        embedder,
		trManager:       trManager,
	}
}
