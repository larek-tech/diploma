package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type HandlerFunc func(context.Context, *kafka.Message) error

type Consumer struct {
	cfg      Config
	consumer *kafka.Consumer
	topic    string
	handler  HandlerFunc
}

func NewConsumer(cfg Config, handler HandlerFunc) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          "dataGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: c,
		topic:    cfg.Topic,
		handler:  handler,
		cfg:      cfg,
	}, nil
}

func (c Consumer) Run(ctx context.Context) error {
	err := c.consumer.SubscribeTopics([]string{c.topic}, nil)

	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}
	defer c.consumer.Close()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("consumer stopped: context canceled: %w", ctx.Err())
		default:
		}
		msg, err := c.consumer.ReadMessage(time.Second)
		if err != nil {
			if !err.(kafka.Error).IsTimeout() {
				slog.Error("Consumer error", "err", err, "msg", msg) // : %v (%v)\n
			}
		}
		if msg == nil {
			continue
		}
		err = c.handler(ctx, msg)
		if err != nil {
			slog.Error("failed to process kafka msg: %w", err)
		}
	}
}
