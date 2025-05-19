package controller

import (
	"context"
	"errors"
	"slices"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/larek-tech/diploma/domain/internal/auth"
	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/source/model"
	"github.com/larek-tech/diploma/domain/pkg/kafka"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	sourceTopic = "source"
	statusTopic = "status"
)

var (
	// ErrUpdateSourceStatus is an error when source status updating failed.
	ErrUpdateSourceStatus = errors.New("failed to update source status while parsing")
	// ErrNoAccessToSource is an error when user can't edit source.
	ErrNoAccessToSource = errors.New("user has no access to edit source")
)

type sourceRepo interface {
	InsertSource(ctx context.Context, s model.SourceDao) (int64, error)
	GetSourceByID(ctx context.Context, id, userID int64, roleIDs []int64) (model.SourceDao, error)
	GetSourceIDs(ctx context.Context, sourceID, userID int64, roleIDs []int64) (uuid.UUID, error)
	UpdateSource(ctx context.Context, s model.SourceDao, userID int64, roleIDs []int64) error
	DeleteSource(ctx context.Context, id, userID int64, roleIDs []int64) error
	ListSources(ctx context.Context, userID int64, roleIDs []int64, offset, limit uint64) ([]model.SourceDao, error)
	ListSourcesByDomain(ctx context.Context, userID, domainID int64, roleIDs []int64, offset, limit uint64) ([]model.SourceDao, error)
	GetPermittedUsers(ctx context.Context, sourceID int64) ([]int64, error)
	GetPermittedRoles(ctx context.Context, sourceID int64) ([]int64, error)
	UpdatePermittedUsers(ctx context.Context, sourceID int64, userIDs []int64) ([]int64, error)
	UpdatePermittedRoles(ctx context.Context, sourceID int64, roleIDs []int64) ([]int64, error)
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

func (ctrl *Controller) checkSourceCreator(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) error {
	userID := meta.GetUserId()

	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.checkSourceCreator",
		trace.WithAttributes(
			attribute.Int64("userID", userID),
			attribute.Int64("sourceID", sourceID),
		),
	)
	defer span.End()

	roles := meta.GetRoles()
	if !slices.Contains(roles, auth.AdminRoleID) {
		source, err := ctrl.sr.GetSourceByID(ctx, sourceID, userID, roles)
		if err != nil {
			return errs.WrapErr(err)
		}

		if source.UserID != userID {
			return errs.WrapErr(ErrNoAccessToSource, "check source creator")
		}
	}
	return nil
}
