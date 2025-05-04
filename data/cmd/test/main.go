package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"

	kafka2 "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/larek-tech/diploma/data/internal/infrastructure/kafka"
)

func main() {
	ctx := context.Background()
	kafkaCfg, err := getKafkaConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	producer, err := kafka.NewProducer(*kafkaCfg)
	if err != nil {
		slog.Error("failed to create kafka producer: %w", err)
		return
	}
	str := "https://gitflic.com"
	content := base64.StdEncoding.EncodeToString([]byte(str))
	testPayload := &DataMessage{
		Title:        "gitflic",
		Content:      []byte(content),
		Type:         Web,
		Credentials:  nil,
		UpdateParams: nil,
	}
	payload, err := json.Marshal(testPayload)
	if err != nil {
		slog.Error("failed to marshal payload: %w", err)
		return
	}

	err = producer.Produce(ctx, &kafka2.Message{
		TopicPartition: kafka2.TopicPartition{
			Topic:       &kafkaCfg.Topic,
			Partition:   0,
			Offset:      0,
			Metadata:    nil,
			Error:       nil,
			LeaderEpoch: nil,
		},
		Value: payload,
		Key:   []byte("data"),
	})
	if err != nil {
		slog.Error("failed to produce message: %w", err)
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
