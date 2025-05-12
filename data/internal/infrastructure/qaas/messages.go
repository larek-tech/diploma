package qaas

import (
	"time"

	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/object_store"
	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type Entity interface {
	site.Site | site.Page | document.Document | object_store.ObjectStore | object_store.Object
}

type DelayedJob[T Entity] struct {
	Payload  *T             `json:"payload"`  // полезная нагрузка
	Delay    int            `json:"delay"`    // задержка в секундах
	Metadata map[string]any `json:"metadata"` // метаданные
}

type SiteJob = DelayedJob[site.Site]
type PageJob = DelayedJob[site.Page]
type PageResultJob = DelayedJob[site.Page]

type ParseS3Job = DelayedJob[object_store.ObjectStore]
type ObjectJob = DelayedJob[object_store.Object]

type EmbedJob = DelayedJob[document.Document]

type ParseStatusJob struct {
	ExternalKey      string // идентификатор полученный от сторонней системы для обработки процесса обработки
	SourceID         string
	SiteJobID        string
	SiteID           string
	ParsePageJobsIDs []string
	Delay            time.Duration
}

type ResultMessage struct {
	SourceID string // uuid ID источника
	ObjID    string // uuid объекта который надо обработать

}
