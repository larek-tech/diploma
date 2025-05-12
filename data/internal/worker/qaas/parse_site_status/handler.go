package parse_site_status

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/larek-tech/diploma/data/internal/infrastructure/ptr"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/infrastructure/queue/messages"
	"go.dataddo.com/pgq"
)

const (
	resultTopic string = "status"
	statusDelay        = time.Second * 10
)

type Handler struct {
	publisher     publisher
	pageJobStore  pageJobStore
	kafkaProducer kafkaProducer
}

func New(publisher publisher, pageJobStore pageJobStore, kafkaProducer kafkaProducer) *Handler {
	return &Handler{
		pageJobStore:  pageJobStore,
		publisher:     publisher,
		kafkaProducer: kafkaProducer,
	}
}

func (h Handler) Handle(ctx context.Context, queueMsg *pgq.MessageIncoming) (bool, error) {
	// Parse the message payload
	var payload qaas.ParseStatusJob
	err := json.Unmarshal(queueMsg.Payload, &payload)
	if err != nil {
		return true, fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	if payload.ExternalKey == "" {
		return true, fmt.Errorf("parse_site_status failed to receive ExternalKey from payload")
	}

	processed, err := h.pageJobStore.GetProcessedPageCount(ctx, payload.SiteJobID)
	if err != nil {
		return true, fmt.Errorf("failed to get processed page count: %w", err)
	}
	unprocessed, err := h.pageJobStore.GetUnprocessedPageCount(ctx, payload.SiteJobID)
	if err != nil {
		return true, fmt.Errorf("failed to get unprocessed page count: %w", err)
	}

	var status messages.SourceStatus
	if processed > 0 && unprocessed == 0 {
		status = messages.StatusReady
	}
	if unprocessed != 0 {
		status = messages.StatusParsing
	}

	key, value, err := h.assembleMessage(payload.ExternalKey, messages.ParsingStatus{
		SourceID:  payload.SourceID,
		Status:    status,
		JobID:     "",
		Processed: processed,
		Total:     processed + unprocessed,
	})
	if err != nil {
		return true, fmt.Errorf("failed to assemble message: %w", err)
	}
	err = h.kafkaProducer.Produce(ctx, resultTopic, key, value)
	if err != nil {
		return false, fmt.Errorf("failed to produce message: %w", err)
	}
	if status == messages.StatusReady {
		return true, nil
	}
	publishOptions := []qaas.PublishOption{
		qaas.WithQueue(qaas.ParseSiteStatusQueue),
		qaas.WithScheduledFor(ptr.To(time.Now().Add(statusDelay))),
	}

	_, err = h.publisher.Publish(ctx, []any{payload}, publishOptions...)
	if err != nil {
		return false, fmt.Errorf("failed to schedule next status check: %w", err)
	}
	return true, nil
}

func (h Handler) assembleMessage(externalKey string, msg messages.ParsingStatus) (key []byte, value []byte, err error) {

	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal payload of ParsingStatus: %w", err)
	}

	return []byte(externalKey), payload, nil
}
