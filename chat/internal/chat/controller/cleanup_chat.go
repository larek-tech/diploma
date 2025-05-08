package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CleanupChat delete empty chat.
func (ctrl *Controller) CleanupChat(ctx context.Context, chatID uuid.UUID) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CleanupChat",
		trace.WithAttributes(
			attribute.String("chatID", chatID.String()),
		),
	)
	defer span.End()

	if err := ctrl.cr.DeleteChat(ctx, chatID); err != nil {
		return errs.WrapErr(err, "cleanup chat")
	}

	return nil
}
