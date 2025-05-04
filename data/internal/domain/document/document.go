package document

import (
	"errors"
	"time"
)

var (
	ErrDocumentNotFound = errors.New("document not found") // ошибка, когда документ не найден
)

type Type string

const (
	TypePage Type = "web.page" // тип документа - веб-страница
	TypeFile Type = "file"     // тип документа - файл
)

type Document struct {
	ID         string         `db:"id"`          // идентификатор документа в векторном хранилище
	SourceID   string         `db:"source_id"`   // идентификатор источника к которому относится документ
	ObjectID   string         `db:"object_id"`   // идентификатор объекта к которому относится документ
	ObjectType Type           `db:"object_type"` // тип объекта к которому относится документ
	Name       string         `db:"name"`        // название документа для файлов или title для html страниц
	Content    string         `db:"content"`     // содержание документа
	Metadata   map[string]any `db:"metadata"`    // метаданные документа (например, заголовок, автор, дата создания и т.д.)
	Chunks     []string       `db:"chunks"`      // IDS чанков данного документа
	CreatedAt  time.Time      `db:"created_at"`  // дата создания документа
	UpdatedAt  time.Time      `db:"updated_at"`  // дата последнего обновления документа
}

type Chunk struct {
	ID         string         `db:"id"`          // идентификатор чанка в векторном хранилище
	Index      int            `db:"index"`       // индекс чанка в документе
	DocumentID string         `db:"document_id"` // идентификатор документа к которому относиться данный чанк
	Content    string         `db:"content"`     // текстовый контент чанка
	Metadata   map[string]any `db:"metadata"`    // метаданные чанка (например, заголовок, автор, дата создания и т.д.)
	Embeddings []float32      `db:"embeddings"`  // векторное представление чанка
}

type Questions struct {
	ID         string    `db:"id"`
	ChunkID    string    `db:"chunk_id"`
	Question   string    `db:"question"`
	Embeddings []float32 `db:"embeddings"`
}
