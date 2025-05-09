package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CleanupChat delete empty chat.
func (ctrl *Controller) CleanupChat(ctx context.Context, chatIDRaw string) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CleanupChat",
		trace.WithAttributes(
			attribute.String("chatID", chatIDRaw),
		),
	)
	defer span.End()

	chatID, err := uuid.Parse(chatIDRaw)
	if err != nil {
		return errs.WrapErr(err, "parse chat id")
	}

	if err := ctrl.cr.DeleteChat(ctx, chatID); err != nil {
		return errs.WrapErr(err, "cleanup chat")
	}

	return nil
}
