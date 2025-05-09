package handler

import (
	"context"

	"github.com/larek-tech/diploma/chat/internal/auth"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateChat creates new chat.
func (h *Handler) CreateChat(ctx context.Context, _ *emptypb.Empty) (*pb.Chat, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.cc.CreateChat(ctx, meta.GetUserId())
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create chat")
		return nil, status.Error(codes.Internal, "failed to create chat")
	}

	return resp, status.Error(codes.OK, "created chat successfully")
}
