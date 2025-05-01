package document

type Document struct {
	ID       string   // идентификатор документа в векторном хранилище
	SourceID string   // идентификатор источника к которому относится документ
	Name     string   // название документа для файлов или title для html страниц
	Content  string   // содержание документа
	Chunks   []string // IDS чанков данного документа
}

type Chunk struct {
	ID         string    // идентификатор чанка в векторном хранилище
	Index      int       // индекс чанка в документе
	DocumentID string    // идентификатор документа к которому относиться данный чанк
	Content    string    // текстовый контент чанка
	Embeddings []float32 // векторное представление чанка
}
