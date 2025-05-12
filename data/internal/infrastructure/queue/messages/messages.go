package messages

import (
	"github.com/larek-tech/diploma/data/internal/domain/source"
)

// отправляем в source_topic
type DataMessage struct {
	Title        string              `json:"title"`
	Content      []byte              `json:"content"` // byte-строка с url или считанный файл
	Type         source.Type         `json:"type"`
	Credentials  []byte              `json:"credentials"`
	UpdateParams source.UpdateParams `json:"update_params"`
}

// отправялем в source_topic
type CreateResponse struct {
	SourceID string
}

type SourceStatus uint8

const (
	// StatusUndefined undefined status.
	StatusUndefined SourceStatus = iota
	// StatusReady source is ready.
	StatusReady
	// StatusParsing source is being parsed.
	StatusParsing
	// StatusFailed source parsing failed.
	StatusFailed
)

// отправялем в status_topic
type ParsingStatus struct {
	SourceID  string       `json:"source"` // uuid ID источника
	Status    SourceStatus `json:"status"` //
	JobID     string       // uuid ID процесса парсинга
	Processed int          // количество элементов обработанных за текущий проход
	Total     int          // количество элементов полученное при первом обходе ресурса
}
