package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/samber/lo"
)

const handlerTimeout = 30 * time.Second

type HandlerFunc func(context.Context, *kafka.Message) error

type Consumer struct {
	cfg      Config
	consumer *kafka.Consumer
	// handler  HandlerFunc
	handlers map[string]HandlerFunc // map of topic to handler
}

func NewConsumer(cfg Config, handlers map[string]HandlerFunc) (*Consumer, error) {
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
		handlers: handlers,
		cfg:      cfg,
	}, nil
}

func (c Consumer) Run(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return fmt.Errorf("no handlers provided")
	}
	err := c.consumer.SubscribeTopics(lo.Keys(c.handlers), nil)

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
		handler, ok := c.handlers[*msg.TopicPartition.Topic]
		if !ok {
			slog.Error("no handler for topic", "topic", *msg.TopicPartition.Topic)
			continue
		}
		ctx, cancel := context.WithTimeout(ctx, handlerTimeout)
		defer cancel()
		err = handler(ctx, msg)
		if err != nil {
			slog.Error("failed to process kafka msg", "err", err)
		}
	}
}
