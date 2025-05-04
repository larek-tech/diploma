package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteDomain deletes domain by id.
func (ctrl *Controller) DeleteDomain(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteDomain",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("domainID", domainID),
		),
	)
	defer span.End()

	if err := ctrl.dr.DeleteDomain(ctx, domainID, meta.GetUserId(), meta.GetRoles()); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
