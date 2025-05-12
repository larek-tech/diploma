package source

import (
	"time"
)

type Type uint8

const (
	Undefined = iota
	Web
	SingleFile
	ArchivedFiles
	S3WithCredentials
)

type UpdateParams struct {
	EveryPeriod int       `json:"every_period"` // обновлять каждые X секунд
	OnTime      time.Time `json:"on_time"`      // обновлять при наступлении даты+времени
}

// отправляем в source_topic
type DataMessage struct {
	ExternalKey  []byte
	Title        string       `json:"title"`
	Content      []byte       `json:"content"` // byte-строка с url или считанный файл
	Type         Type         `json:"type"`
	Credentials  []byte       `json:"credentials"`
	UpdateParams UpdateParams `json:"update_params"`
}

type Source struct {
	ID          string `db:"id"`          // ID uuid идентификатор источника
	Title       string `db:"title"`       // Title название источника
	Type        Type   `db:"type"`        // Type тип источника (с паролем, без пароля, архив)
	Credentials []byte `db:"credentials"` // Credentials учетные данные для доступа к источнику
}
