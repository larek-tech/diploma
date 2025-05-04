package controller

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/larek-tech/diploma/domain/internal/domain/source/model"
	"github.com/larek-tech/diploma/domain/pkg/kafka"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/trace"
)

const (
	sourceTopic = "source"
	statusTopic = "status"
)

var (
	// ErrUpdateSourceStatus is an error when source status updating failed.
	ErrUpdateSourceStatus = errors.New("failed to update source status while parsing")
)

type sourceRepo interface {
	InsertSource(ctx context.Context, s model.SourceDao) (int64, error)
	GetSourceByID(ctx context.Context, id, userID int64, roleIDs []int64) (model.SourceDao, error)
	UpdateSource(ctx context.Context, s model.SourceDao, userID int64, roleIDs []int64) error
	DeleteSource(ctx context.Context, id, userID int64, roleIDs []int64) error
	ListSources(ctx context.Context, userID int64, roleIDs []int64, offset, limit uint64) ([]model.SourceDao, error)
}

// Controller implements source methods on logic layer.
type Controller struct {
	sr       sourceRepo
	tracer   trace.Tracer
	producer *kafka.AsyncProducer
	consumer *kafka.Consumer
	statusCh chan *sarama.ConsumerMessage
	errCh    chan error
}

// New creates new Controller.
func New(ctx context.Context, sr sourceRepo, tracer trace.Tracer, producer *kafka.AsyncProducer, consumer *kafka.Consumer) (*Controller, error) {
	statusCh, errCh, err := consumer.Subscribe(ctx, statusTopic)
	if err != nil {
		return nil, errs.WrapErr(err, "subscribe to status topic")
	}

	return &Controller{
		sr:       sr,
		tracer:   tracer,
		producer: producer,
		consumer: consumer,
		statusCh: statusCh,
		errCh:    errCh,
	}, nil
}
