package server

import "time"

type SourceType uint8

const (
	SourceWithCredentials = iota
	SourceSingleFile
	SourceArchivedFiles
)

type UpdateParams struct {
	EveryPeriod int       `json:"every_period"` // обновлять каждые X секунд
	OnTime      time.Time `json:"on_time"`      // обновлять при наступлении даты+времени
}

type DataMessage struct {
	Title        string       `json:"title"`
	Content      []byte       `json:"content"` // byte-строка с url или считанный файл
	Type         SourceType   `json:"type"`
	Credentials  []byte       `json:"credentials"`
	UpdateParams UpdateParams `json:"update_params"`
}
