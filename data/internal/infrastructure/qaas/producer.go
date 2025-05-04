package qaas

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
	"go.dataddo.com/pgq"
)

type Publisher struct {
	pub pgq.Publisher
}

func NewPublisher(pub pgq.Publisher) *Publisher {
	return &Publisher{
		pub: pub,
	}
}

func (p *Publisher) Publish(ctx context.Context, msg any, scheduledAt ...*time.Time) error {
	var rawMsg json.RawMessage
	var msgType string

	switch m := msg.(type) {
	case messages.PageJob:
		data, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to marshal page job: %w", err)
		}
		rawMsg = data
		msgType = string(messages.ParsePage)
	case messages.SiteJob:
		data, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to marshal site job: %w", err)
		}
		rawMsg = data
		msgType = string(messages.ParseSite)
	case messages.ResultMessage:
		data, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to marshal result message: %w", err)
		}
		rawMsg = data
		msgType = string(m.Type)
	}
	var ScheduledFor *time.Time
	if len(scheduledAt) > 0 {
		ScheduledFor = scheduledAt[0]
	}
	_, err := p.pub.Publish(ctx, QueueName, &pgq.MessageOutgoing{
		ScheduledFor: ScheduledFor,
		Payload:      rawMsg,
		Metadata: pgq.Metadata{
			"type": msgType,
		},
	})
	return err
}
