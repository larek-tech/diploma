package qaas

import (
	"context"
	"database/sql"
	"fmt"

	"go.dataddo.com/pgq"
)

type Consumer struct {
	db *sql.DB
}

func NewConsumer(
	db *sql.DB,
) *Consumer {
	return &Consumer{
		db: db,
	}
}

type adapter struct {
	h handler
}

func (m adapter) HandleMessage(ctx context.Context, msg *pgq.MessageIncoming) (processed bool, err error) {
	return m.h.Handle(ctx, msg)
}

func (c *Consumer) Run(ctx context.Context, queue Queue, h handler) error {
	consumer, err := pgq.NewConsumer(c.db, string(queue), adapter{h: h})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	return consumer.Run(ctx)
}

// TODO: do we embed message type if we know the queue?
// func (c *Consumer) HandleMessage(ctx context.Context, msg *pgq.MessageIncoming) (processed bool, err error) {
// 	msgQueue, ok := msg.Metadata["queue"]
// 	if !ok {
// 		return false, nil
// 	}
// 	switch Queue(msgQueue) {
// 	// case ParsePageQueue:
// 	// 	var job PageJob
// 	// 	err = json.Unmarshal(msg.Payload, &job)
// 	// 	if err != nil {
// 	// 		return true, fmt.Errorf("failed to unmarshal parsepage payload: %w", err)
// 	// 	}
// 	// 	if err = c.pageJobHandler.Handle(ctx, job); err != nil {
// 	// 		return true, fmt.Errorf("failed to handle page job: %w", err)
// 	// 	}
// 	// case ParseSiteQueue:
// 	// 	var job SiteJob
// 	// 	err = json.Unmarshal(msg.Payload, &job)
// 	// 	if err != nil {
// 	// 		return true, fmt.Errorf("failed to unmarshal parsesite payload: %w", err)
// 	// 	}
// 	// 	if err = c.siteJobHandler.Handle(ctx, job); err != nil {
// 	// 		return true, fmt.Errorf("failed to handle site job: %w", err)
// 	// 	}
// 	case ParsePageResultQueue:
// 		// var result PageResultJob
// 		// err = json.Unmarshal(msg.Payload, &result)
// 		// if err != nil {
// 		// 	return true, fmt.Errorf("failed to unmarshal result message payload: %w", err)
// 		// }
// 		// if err = c.resultMessageHandler.Handle(ctx, result); err != nil {
// 		// 	return true, fmt.Errorf("failed to handle result message: %w", err)
// 		// }

// 	default:
// 		return true, fmt.Errorf("unknown message type: %s", msgQueue)
// 	}
// 	return true, nil
// }
