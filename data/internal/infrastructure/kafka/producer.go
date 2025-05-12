package kafka

import (
	"context"
	"fmt"
	"log/slog" // Added for logging

	"github.com/IBM/sarama" // Changed import
)

// Producer struct updated for Sarama (using SyncProducer)
type Producer struct {
	syncProducer sarama.SyncProducer
}

// NewProducer updated for Sarama
func NewProducer(cfg Config) (*Producer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true // Required for SyncProducer to wait for ack
	saramaCfg.Producer.Return.Errors = true
	// saramaCfg.Producer.RequiredAcks = sarama.WaitForAll // Example: wait for all in-sync replicas
	// saramaCfg.Producer.Retry.Max = 5                   // Example: configure retries
	// It's good practice to set the Kafka version. Choose a version compatible with your cluster.
	// Example: saramaCfg.Version = sarama.V2_8_1_0

	producer, err := sarama.NewSyncProducer([]string{cfg.BootstrapServers}, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init sarama sync producer: %w", err)
	}

	return &Producer{
		syncProducer: producer,
	}, nil
}

// Produce method adapted for Sarama.
// The signature is changed to accept topic, key, and value for simplicity.
// The original context `ctx` is used here to check for cancellation before the blocking send.
func (p *Producer) Produce(ctx context.Context, topic string, key []byte, value []byte) error {
	// Check context before attempting to send, as SendMessage is blocking
	select {
	case <-ctx.Done():
		return fmt.Errorf("producing message cancelled: %w", ctx.Err())
	default:
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),   // Use sarama.StringEncoder if key is string
		Value: sarama.ByteEncoder(value), // Use sarama.StringEncoder if value is string
	}

	// SendMessage is synchronous for SyncProducer
	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka topic %s: %w", topic, err)
	}

	slog.Debug("Message sent successfully to Kafka", "topic", topic, "partition", partition, "offset", offset)
	return nil
}

// ProduceMessage allows sending a pre-constructed sarama.ProducerMessage.
// This is useful if you need more control (e.g., headers, specific partition, timestamp).
func (p *Producer) ProduceMessage(ctx context.Context, msg *sarama.ProducerMessage) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("producing sarama message cancelled: %w", ctx.Err())
	default:
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send sarama message to topic %s: %w", msg.Topic, err)
	}
	slog.Debug("Sarama message sent successfully", "topic", msg.Topic, "partition", partition, "offset", offset)
	return nil
}

// Close method for Sarama producer
func (p *Producer) Close() {
	if p.syncProducer != nil {
		slog.Info("Closing sarama sync producer")
		if err := p.syncProducer.Close(); err != nil {
			slog.Error("Failed to close sarama sync producer", "err", err)
		}
	}
}
