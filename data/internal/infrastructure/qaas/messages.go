package qaas

import (
	"github.com/larek-tech/diploma/data/internal/domain/document"
	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type Entity interface {
	site.Site | site.Page | document.Document
}

type DelayedJob[T Entity] struct {
	Payload  *T             `json:"payload"`  // полезная нагрузка
	Delay    int            `json:"delay"`    // задержка в секундах
	Metadata map[string]any `json:"metadata"` // метаданные
}

type SiteJob = DelayedJob[site.Site]
type PageJob = DelayedJob[site.Page]
type PageResultJob = DelayedJob[site.Page]
type EmbedJob = DelayedJob[document.Document]

type ResultMessage struct {
	SourceID string // uuid ID источника
	ObjID    string // uuid объекта который надо обработать

}
