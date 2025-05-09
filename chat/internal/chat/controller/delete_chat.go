package controller

import (
	"context"

	"github.com/google/uuid"
	authpb "github.com/larek-tech/diploma/chat/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// DeleteChat soft delete chat.
func (ctrl *Controller) DeleteChat(ctx context.Context, chatIDRaw string, meta *authpb.UserAuthMetadata) error {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.DeleteChat",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("chatID", chatIDRaw),
		),
	)
	defer span.End()

	chatID, err := uuid.Parse(chatIDRaw)
	if err != nil {
		return errs.WrapErr(err, "parse chat id")
	}

	chatCreatorID, err := ctrl.cr.GetChatUserID(ctx, chatID)
	if err != nil {
		return errs.WrapErr(err)
	}

	if chatCreatorID != meta.GetUserId() {
		return errs.WrapErr(ErrNoAccessToChat)
	}

	if err = ctrl.cr.SoftDeleteChat(ctx, chatID); err != nil {
		return errs.WrapErr(err)
	}

	return nil
}
