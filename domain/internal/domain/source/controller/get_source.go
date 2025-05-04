package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetSource returns source by id.
func (ctrl *Controller) GetSource(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) (*pb.GetSourceResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetSource",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("sourceID", sourceID),
		),
	)
	defer span.End()

	source, err := ctrl.sr.GetSourceByID(ctx, sourceID, meta.GetUserId(), meta.GetRoles())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return &pb.GetSourceResponse{Source: source.ToProto()}, nil
}
