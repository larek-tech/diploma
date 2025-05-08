package controller

import (
	"context"
	"time"

	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CreateChat creates new chat.
func (ctrl *Controller) CreateChat(ctx context.Context, userID int64) (*pb.Chat, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CreateChat",
		trace.WithAttributes(attribute.Int64("userID", userID)),
	)
	defer span.End()

	chat := model.ChatDao{
		UserID:    userID,
		Title:     model.ChatDefaultTitle,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	chatID, err := ctrl.cr.InsertChat(ctx, chat)
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	chat.ID = chatID
	return chat.ToProto(), nil
}
