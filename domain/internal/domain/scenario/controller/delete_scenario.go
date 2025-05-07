package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteScenario deletes scenario by id.
func (ctrl *Controller) DeleteScenario(ctx context.Context, scenarioID int64, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteScenario",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("scenarioID", scenarioID),
		),
	)
	defer span.End()

	if err := ctrl.sr.DeleteScenario(ctx, scenarioID, meta.GetUserId()); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
