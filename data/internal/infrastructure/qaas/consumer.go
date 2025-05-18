package qaas

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.dataddo.com/pgq"
)

const (
	parallelMessages = 10
	pollingInterval  = 1 * time.Second
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
	consumer, err := pgq.NewConsumer(c.db, string(queue), adapter{h: h},
		pgq.WithMaxParallelMessages(parallelMessages),
		pgq.WithPollingInterval(pollingInterval),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	return consumer.Run(ctx)
}
