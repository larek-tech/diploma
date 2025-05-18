package create_source

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/queue/messages"
	"github.com/larek-tech/diploma/data/pkg/metric"
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

func (h Handler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	slog.Info("received new msg", "msg", string(msg.Key))

	var payload source.DataMessage
	if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&payload); err != nil {
		return fmt.Errorf("failed to decode json")
	}
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		return fmt.Errorf("failed to process DataMessage: %w", err)
	}
	payload.ExternalKey = msg.Key

	newSource, err := h.service.CreateSource(ctx, payload)
	metric.IncrementSourcesCreated("undefined", newSource.ID, err)
	if err != nil {
		return fmt.Errorf("failed to create new source: %w", err)
	}
	key, value, err := h.assembleMessage(msg, newSource)
	if err != nil {
		return fmt.Errorf("failed to create response msg: %w", err)
	}
	err = h.kafkaProducer.Produce(ctx, resultTopic, key, value)
	if err != nil {
		return fmt.Errorf("create_source failed to send status msg: %w", err)
	}

	return nil
}

func (h Handler) assembleMessage(incomingMsg *sarama.ConsumerMessage, src *source.Source) (key []byte, value []byte, err error) {
	if src == nil {
		return nil, nil, fmt.Errorf("failed to assemble create source message source is nil")
	}
	msg := messages.ParsingStatus{
		SourceID:  src.ID,
		Status:    messages.StatusParsing,
		JobID:     "",
		Processed: 0,
		Total:     0,
	}
	value, err = json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal payload of ParsingStatus: %w", err)
	}
	key = incomingMsg.Key
	return
}
