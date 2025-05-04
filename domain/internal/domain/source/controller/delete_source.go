package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteSource deletes source by id.
func (ctrl *Controller) DeleteSource(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteSource",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("sourceID", sourceID),
		),
	)
	defer span.End()

	if err := ctrl.sr.DeleteSource(ctx, sourceID, meta.GetUserId(), meta.GetRoles()); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
