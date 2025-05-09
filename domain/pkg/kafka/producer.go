package kafka

import (
	"errors"
	"net"
	"strconv"

	"github.com/IBM/sarama"
)

var (
	// ErrAsyncProducer is an error when Kafka async producer wasn't opened.
	ErrAsyncProducer = errors.New("creating new async Kafka producer failed")
	// ErrCloseProducer is an error when async producer wasn't closed properly.
	ErrCloseProducer = errors.New("can't close async producer properly")
)

// AsyncProducer is a Kafka async producer.
type AsyncProducer struct {
	Config   *Config
	producer sarama.AsyncProducer
}

// NewAsyncProducer creates a new Kafka async producer.
func NewAsyncProducer(config *Config, partitioner sarama.PartitionerConstructor, acks sarama.RequiredAcks) (*AsyncProducer, chan error, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = partitioner
	cfg.Producer.RequiredAcks = acks

	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true

	brokers := make([]string, len(config.Brokers))
	for idx, broker := range config.Brokers {
		brokers[idx] = net.JoinHostPort(broker.Host, strconv.Itoa(broker.Port))
	}

	asyncProducer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, nil, errors.Join(ErrAsyncProducer, err)
	}

	errCh := make(chan error)
	go func() {
		for e := range asyncProducer.Errors() {
			errCh <- e
		}
	}()

	return &AsyncProducer{
		Config:   config,
		producer: asyncProducer,
	}, errCh, nil
}

// SendAsyncMessage sends a message to Kafka.
func (k *AsyncProducer) SendAsyncMessage(message *sarama.ProducerMessage) {
	k.producer.Input() <- message
}

// Close closes the Kafka async producer.
func (k *AsyncProducer) Close() error {
	if err := k.producer.Close(); err != nil {
		return errors.Join(ErrCloseProducer, err)
	}
	return nil
}
