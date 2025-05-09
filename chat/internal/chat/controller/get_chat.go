package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetChat returns chat and its content.
func (ctrl *Controller) GetChat(ctx context.Context, chatIDRaw string) (*pb.Chat, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetChat",
		trace.WithAttributes(attribute.String("chatID", chatIDRaw)),
	)
	defer span.End()

	chatID, err := uuid.Parse(chatIDRaw)
	if err != nil {
		return nil, errs.WrapErr(err, "parse chat id")
	}

	chat, err := ctrl.cr.GetChat(ctx, chatID)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return chat.ToProto(), nil
}
