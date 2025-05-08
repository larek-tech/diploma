package controller

import (
	"context"

	"github.com/google/uuid"
	authpb "github.com/larek-tech/diploma/chat/internal/auth/pb"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RenameChat change chat title.
func (ctrl *Controller) RenameChat(ctx context.Context, req *pb.RenameChatRequest, meta *authpb.UserAuthMetadata) (*pb.Chat, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.RenameChat",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("chatID", req.GetChatId()),
			attribute.String("title", req.GetTitle()),
		),
	)
	defer span.End()

	chatID, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, errs.WrapErr(err, "invalid chat id")
	}

	chat, err := ctrl.cr.GetChat(ctx, chatID)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	if chat.UserID != meta.GetUserId() {
		return nil, errs.WrapErr(ErrNoAccessToChat)
	}

	if err = ctrl.cr.UpdateChatTitle(ctx, req.GetTitle(), chatID); err != nil {
		return nil, errs.WrapErr(err)
	}

	return chat.ToProto(), nil
}
