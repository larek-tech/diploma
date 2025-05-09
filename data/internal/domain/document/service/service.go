package service

type Service struct {
	documentStorage documentStorage
	chunkStorage    chunkStorage
	questionStorage questionStorage
	embedder        embedder
	llm             llm
	trManager       trManager
}

func New(documentStorage documentStorage, chunkStorage chunkStorage, questionStorage questionStorage, embedder embedder, llm llm, trManager trManager) *Service {
	return &Service{
		documentStorage: documentStorage,
		chunkStorage:    chunkStorage,
		questionStorage: questionStorage,
		embedder:        embedder,
		llm:             llm,
		trManager:       trManager,
	}
}
