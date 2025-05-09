package qaas

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
	"go.dataddo.com/pgq"
)

const QueueName = "qaas"

type Consumer struct {
	pageJobHandler       pageJobHandler
	siteJobHandler       siteJobHandler
	resultMessageHandler resultMessageHandler
	db                   *sql.DB
}

func NewConsumer(
	pageJobHandler pageJobHandler,
	siteJobHandler siteJobHandler,
	resultMessageHandler resultMessageHandler,
	db *sql.DB,
) *Consumer {
	return &Consumer{
		pageJobHandler:       pageJobHandler,
		siteJobHandler:       siteJobHandler,
		resultMessageHandler: resultMessageHandler,
		db:                   db,
	}
}

type Message interface {
	messages.PageJob | messages.SiteJob | messages.ResultMessage
}

func (c *Consumer) HandleMessage(ctx context.Context, msg *pgq.MessageIncoming) (processed bool, err error) {
	msgType, ok := msg.Metadata["type"]
	if !ok {
		return false, nil
	}
	switch messages.MessageType(msgType) {
	case messages.ParsePage:
		var job messages.PageJob
		err = json.Unmarshal(msg.Payload, &job)
		if err != nil {
			return true, fmt.Errorf("failed to unmarshal parsepage payload: %w", err)
		}
		if err = c.pageJobHandler.Handle(ctx, job); err != nil {
			return true, fmt.Errorf("failed to handle page job: %w", err)
		}
	case messages.ParseSite:
		var job messages.SiteJob
		err = json.Unmarshal(msg.Payload, &job)
		if err != nil {
			return true, fmt.Errorf("failed to unmarshal parsesite payload: %w", err)
		}
		if err = c.siteJobHandler.Handle(ctx, job); err != nil {
			return true, fmt.Errorf("failed to handle site job: %w", err)
		}
	case messages.WebResult:
		var result messages.ResultMessage
		err = json.Unmarshal(msg.Payload, &result)
		if err != nil {
			return true, fmt.Errorf("failed to unmarshal result message payload: %w", err)
		}
		if err = c.resultMessageHandler.Handle(ctx, result); err != nil {
			return true, fmt.Errorf("failed to handle result message: %w", err)
		}

	default:
		return true, fmt.Errorf("unknown message type: %s", msgType)
	}
	return true, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	consumer, err := pgq.NewConsumer(c.db, QueueName, c)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	return consumer.Run(ctx)
}
