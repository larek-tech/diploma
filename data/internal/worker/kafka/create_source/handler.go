package create_source

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/queue/messages"
)

const (
	resultTopic string = "status"
	statusDelay        = time.Second * 10
)

type Handler struct {
	service       service
	kafkaProducer kafkaProducer
}

func New(service service, kafkaProducer kafkaProducer) *Handler {
	return &Handler{
		service:       service,
		kafkaProducer: kafkaProducer,
	}
}

func (h Handler) Handle(ctx context.Context, msg *kafka.Message) error {
	slog.Info("received new msg", "msg", msg)
	var payload source.DataMessage
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		return fmt.Errorf("failed to process DataMessage: %w", err)
	}

	newSource, err := h.service.CreateSource(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to create new source: %w", err)
	}
	statusMsg, err := h.assembleMessage(msg, newSource)
	if err != nil {
		return fmt.Errorf("failed to create response msg: %w", err)
	}
	err = h.kafkaProducer.Produce(ctx, statusMsg)
	if err != nil {
		return fmt.Errorf("create_source failed to send status msg: %w", err)
	}

	return nil
}

func (h Handler) assembleMessage(incomingMsg *kafka.Message, src *source.Source) (*kafka.Message, error) {
	if src == nil {
		return nil, fmt.Errorf("failed to assemble create source message source is nil")
	}
	msg := messages.ParsingStatus{
		SourceID:  src.ID,
		JobID:     "",
		Processed: 0,
		Total:     0,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload of ParsingStatus: %w", err)
	}

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{},
		Value:          payload,
		Key:            incomingMsg.Key,
	}, nil
}
