package handler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/chat/internal/auth"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetChat returns chat and its content.
func (h *Handler) GetChat(ctx context.Context, req *pb.GetChatRequest) (*pb.Chat, error) {
	_, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.cc.GetChat(ctx, req.GetChatId())
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get chat")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}
		return nil, status.Error(codes.Internal, "failed to get chat")
	}

	return resp, status.Error(codes.OK, "got chat successfully")
}
