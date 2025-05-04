package handler

import (
	"context"
	"errors"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrUpdateParamsPositive is an error when provided non-positive update params.
	ErrUpdateParamsPositive = errors.New("update period must be positive")
	// ErrCronUpdateParamsNonNegative is an error when provided negative cron-format update params.
	ErrCronUpdateParamsNonNegative = errors.New("cron update params must be non-negative")
)

type sourceController interface {
	CreateSource(ctx context.Context, req *pb.CreateSourceRequest, meta *authpb.UserAuthMetadata) (*pb.Source, error)
	GetSource(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) (*pb.GetSourceResponse, error)
	UpdateSource(ctx context.Context, req *pb.UpdateSourceRequest, meta *authpb.UserAuthMetadata) (*pb.Source, error)
	DeleteSource(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) error
	ListSources(ctx context.Context, req *pb.ListSourcesRequest, meta *authpb.UserAuthMetadata) (*pb.ListSourcesResponse, error)
}

// Handler implements source methods on transport level.
type Handler struct {
	pb.UnimplementedSourceServiceServer
	sc     sourceController
	tracer trace.Tracer
}

// New creates new Handler.
func New(sc sourceController, tracer trace.Tracer) *Handler {
	return &Handler{
		sc:     sc,
		tracer: tracer,
	}
}

func validateUpdateParams(updateParams *pb.UpdateParams) error {
	if updateParams != nil {
		switch {
		case updateParams.EveryPeriod != nil && *updateParams.EveryPeriod <= 0:
			return errs.WrapErr(ErrUpdateParamsPositive)
		case updateParams.GetCron() != nil:
			cron := updateParams.GetCron()
			if cron.GetDayOfWeek() < 0 || cron.Month < 0 || cron.DayOfMonth < 0 || cron.Hour < 0 || cron.Minute < 0 {
				return errs.WrapErr(ErrCronUpdateParamsNonNegative)
			}
		}
	}
	return nil
}
