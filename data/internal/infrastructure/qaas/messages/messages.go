package messages

import (
	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type MessageType string

const (
	ParseSite MessageType = "web.parse_site"  // сообщение для запуска парсинга сайта
	ParsePage MessageType = "web.parse_page"  // сообщение для запуска парсинга страницы
	EmbedPage MessageType = "qaas.embed_page" // сообщение для векторизации страницы

	WebResult MessageType = "web.page.result"

	FileResult  MessageType = "s3.file.result"    // сообщение с результатами парсинга файла
	EmbedResult MessageType = "qaas.embed.result" // сообщение с результатами векторизации
)

type Entity interface {
	site.Site | site.Page
}

type DelayedJob[T Entity] struct {
	Type    MessageType `json:"type"`    // тип сообщения
	Payload *T          `json:"payload"` // полезная нагрузка
	Delay   int         `json:"delay"`   // задержка в секундах
}

type SiteJob = DelayedJob[site.Site]
type PageJob = DelayedJob[site.Page]

type ResultMessage struct {
	SourceID string // uuid ID источника
	Type     MessageType
	ObjID    string // uuid объекта который надо обработать
}
