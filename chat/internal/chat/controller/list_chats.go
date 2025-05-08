package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/chat/internal/auth/pb"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListChats returns paginated list of chats.
func (ctrl *Controller) ListChats(ctx context.Context, req *pb.ListChatsRequest, meta *authpb.UserAuthMetadata) (*pb.ListChatsResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListChats",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("offset", int64(req.GetOffset())),
			attribute.Int64("limit", int64(req.GetLimit())),
		),
	)
	defer span.End()

	chatsDao, err := ctrl.cr.ListChats(ctx, req.GetOffset(), req.GetLimit(), meta.GetUserId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	chats := make([]*pb.Chat, len(chatsDao))
	for idx := range chatsDao {
		chats[idx] = chatsDao[idx].ToProto()
	}

	return &pb.ListChatsResponse{Chats: chats}, nil
}
