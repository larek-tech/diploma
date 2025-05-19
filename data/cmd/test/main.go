package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/larek-tech/diploma/data/internal/infrastructure/kafka"
)

func main() {
	msgType := flag.String("type", "site", "Type of message: site or png")
	flag.Parse()

	var testPayload *DataMessage
	switch *msgType {
	case "site":
		testPayload = getSiteMsg()
	case "png":
		testPayload = getPngMsg()
	default:
		slog.Error("unknown message type", "type", *msgType)
		return
	}

	ctx := context.Background()
	kafkaCfg, err := getKafkaConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	producer, err := kafka.NewProducer(*kafkaCfg)
	if err != nil {
		slog.Error("failed to create kafka producer", "err", err)
		return
	}

	payload, err := json.Marshal(testPayload)
	if err != nil {
		slog.Error("failed to marshal payload", "err", err)
		return
	}

	err = producer.Produce(ctx, "source", []byte("some key"), payload)
	if err != nil {
		slog.Error("failed to produce message", "err", err)
		return
	}
}
func getKafkaConfig() (*kafka.Config, error) {
	var cfg kafka.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read kafka config: %w", err)
	}
	return &cfg, nil
}

func getSiteMsg() *DataMessage {
	str := "https://notes.kiriha.ru/sitemap.xml"
	return &DataMessage{
		Title:        "gitflic",
		Content:      []byte(str),
		Type:         Web,
		Credentials:  nil,
		UpdateParams: nil,
	}
}

func getPngMsg() *DataMessage {
	filePath := "mocks/image.png"
	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("failed to read png file", "err", err)
		return nil
	}
	return &DataMessage{
		Title:        "image.png",
		Content:      data,
		Type:         SingleFile,
		Credentials:  nil,
		UpdateParams: nil,
	}
}

type SourceType uint8

const (
	Undefined SourceType = iota
	Web
	SingleFile
	ArchivedFiles
	WithCredentials
)

// Cron contains cron-format parameters for source updates.
type Cron struct {
	WeekDay int32 `db:"cron_week_day"`
	Month   int32 `db:"cron_month"`
	Day     int32 `db:"cron_day"`
	Hour    int32 `db:"cron_hour"`
	Minute  int32 `db:"cron_minute"`
}

// UpdateParams sets time conditions to parse dynamic source (not static files).
type UpdateParams struct {
	EveryPeriod *int64 `json:"every_period,omitempty"` // update every N seconds
	Cron        *Cron  `json:"cron,omitempty"`         // update on date/time (cron-format)
}

// DataMessage contains information about new Source and is sent to Data service to be processed.
type DataMessage struct {
	Title        string        `json:"title"`
	Content      []byte        `json:"content"` // byte encoded url or file content
	Type         SourceType    `json:"type"`
	Credentials  []byte        `json:"credentials,omitempty"`
	UpdateParams *UpdateParams `json:"update_params,omitempty"`
}
