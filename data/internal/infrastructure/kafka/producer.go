package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(cfg Config) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
	}, nil
}

func (p Producer) Produce(ctx context.Context, msg *kafka.Message) error {
	deliveryChan := make(chan kafka.Event, 1)
	defer close(deliveryChan)

	if err := p.producer.Produce(msg, deliveryChan); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case ev := <-deliveryChan:
		m, ok := ev.(*kafka.Message)
		if !ok {
			return fmt.Errorf("unexpected event type: %T", ev)
		}
		if m.TopicPartition.Error != nil {
			return m.TopicPartition.Error
		}
		return nil
	}
}
func (p Producer) Close() {
	p.producer.Close()
}
