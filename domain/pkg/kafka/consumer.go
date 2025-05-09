package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

var (
	// ErrNewConsumer is an error when new Kafka consumer can't be created.
	ErrNewConsumer = errors.New("creating new consumer failed")
	// ErrGetPartitions is returned when failed to get partitions.
	ErrGetPartitions = errors.New("failed to get partitions")
	// ErrConsumePartition is returned when failed to consume partition.
	ErrConsumePartition = errors.New("error consuming partition")
	// ErrClosePartitionConsumer is an error when a partition consumer wasn't properly closed.
	ErrClosePartitionConsumer = errors.New("partition consumer wasn't properly closed")
)

// Consumer is a Kafka consumer.
type Consumer struct {
	Config         *Config
	SingleConsumer sarama.Consumer
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(config *Config) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = false
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	if config.CommitInterval != 0 {
		cfg.Consumer.Offsets.AutoCommit.Interval = time.Second * time.Duration(config.CommitInterval)
	}

	if config.OffsetNewest {
		cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	brokers := make([]string, len(config.Brokers))
	for idx, broker := range config.Brokers {
		brokers[idx] = net.JoinHostPort(broker.Host, strconv.Itoa(broker.Port))
	}

	consumer, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		return nil, errors.Join(ErrNewConsumer, err)
	}

	return &Consumer{
		Config:         config,
		SingleConsumer: consumer,
	}, nil
}

// Subscribe subscribes to a Kafka topic and sends messages to the out channel in a separate goroutine.
func (consumer *Consumer) Subscribe(ctx context.Context, topic string) (
	out chan *sarama.ConsumerMessage,
	errCh chan error,
	e error,
) {
	partitions, err := consumer.SingleConsumer.Partitions(topic)
	if err != nil {
		return nil, nil, errors.Join(ErrGetPartitions, err)
	}

	initialOffset := sarama.OffsetOldest
	if consumer.Config.OffsetNewest {
		initialOffset = sarama.OffsetNewest
	}

	out = make(chan *sarama.ConsumerMessage)
	errCh = make(chan error)

	for _, partition := range partitions {
		pc, err := consumer.SingleConsumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			close(out)
			close(errCh)

			return nil, nil, fmt.Errorf("partition: %d, %w", partition, errors.Join(ErrConsumePartition, err))
		}

		go consume(ctx, pc, out, errCh)
	}
	return out, errCh, nil
}

func consume(ctx context.Context, pc sarama.PartitionConsumer, out chan<- *sarama.ConsumerMessage, err chan<- error) {
	for {
		select {
		case <-ctx.Done():
			if e := pc.Close(); e != nil {
				err <- errors.Join(ErrClosePartitionConsumer, e)
			}
			return
		case message := <-pc.Messages():
			out <- message
		}
	}
}
