package parse_site_status

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/larek-tech/diploma/data/internal/infrastructure/ptr"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/infrastructure/queue/messages"
	"go.dataddo.com/pgq"
)

const (
	resultTopic string = "status_topic"
	statusDelay        = time.Second * 10
)

type Handler struct {
	publisher     publisher
	pageJobStore  pageJobStore
	kafkaProducer kafkaProducer
}

func New(publisher publisher, pageJobStore pageJobStore, kafkaProducer kafkaProducer) *Handler {
	return &Handler{
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
	// Check if the job is already processed
	return true, nil
	// processed, err := h.pageJobStore.GetProcessedPageCount(ctx, payload.SiteID)
	// if err != nil {
	// 	return true, fmt.Errorf("failed to get processed page count: %w", err)
	// }
	// unprocessed, err := h.pageJobStore.GetUnprocessedPageCount(ctx, payload.SiteID)
	// if err != nil {
	// 	return true, fmt.Errorf("failed to get unprocessed page count: %w", err)
	// }

	// responseMsg, err := h.assembleMessage(payload.SiteID, processed, unprocessed)
	// if err != nil {
	// 	return true, fmt.Errorf("failed to assemble message: %w", err)
	// }
	// err = h.kafkaProducer.Produce(ctx, responseMsg)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to produce message: %w", err)
	// }
	// if processed > 0 && unprocessed == 0 {
	// 	return true, nil
	// }
	// // schedule next status check
	// publishOptions := []qaas.PublishOption{
	// 	qaas.WithQueue(qaas.ParseSiteStatusQueue),
	// 	qaas.WithScheduledFor(ptr.To(time.Now().Add(statusDelay))),
	// }

	// _, err = h.publisher.Publish(ctx, []any{payload}, publishOptions...)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to schedule next status check: %w", err)
	// }

	// return true, nil
}

func (h Handler) assembleMessage(siteID string, processed, unprocessed int) (*kafka.Message, error) {
	msg := messages.ParsingStatus{
		SourceID:  siteID,
		JobID:     "",
		Processed: processed,
		Total:     processed + unprocessed,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload of ParsingStatus: %w", err)
	}

	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     ptr.To(resultTopic),
			Partition: kafka.PartitionAny,
		},
		Value: payload,
		Key:   []byte("data"),
	}, nil
}
