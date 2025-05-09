package handler

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/auth"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListChats returns paginated list of chats.
func (h *Handler) ListChats(ctx context.Context, req *pb.ListChatsRequest) (*pb.ListChatsResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.cc.ListChats(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("list chats")
		return nil, status.Error(codes.Internal, "failed to list chats")
	}

	return resp, status.Error(codes.OK, "returned chats list successfully")
}
