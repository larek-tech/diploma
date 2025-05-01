package messages

import (
	"github.com/larek-tech/diploma/data/internal/domain/site"
)

type MessageType string

const (
	ParseSite MessageType = "web.parse_site" // сообщение для запуска парсинга сайта
	ParsePage MessageType = "web.parse_page" // сообщение для запуска парсинга страницы
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
	SourceID  string // uuid ID источника
	JobID     string // uuid ID процесса парсинга
	Processed int    // количество элементов обработанных за текущий проход
	Total     int    // количество элементов полученное при первом обходе ресурса
}
