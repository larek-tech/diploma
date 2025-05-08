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
func (ctrl *Controller) GetChat(ctx context.Context, chatID uuid.UUID) (*pb.Chat, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetChat",
		trace.WithAttributes(attribute.String("chatID", chatID.String())),
	)
	defer span.End()

	chat, err := ctrl.cr.GetChat(ctx, chatID)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return chat.ToProto(), nil
}
