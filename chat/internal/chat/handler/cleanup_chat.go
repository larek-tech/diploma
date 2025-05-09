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
	"google.golang.org/protobuf/types/known/emptypb"
)

// CleanupChat deletes chat if it is empty. Must be called only manually.
func (h *Handler) CleanupChat(ctx context.Context, req *pb.CleanupChatRequest) (*emptypb.Empty, error) {
	_, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if err = h.cc.CleanupChat(ctx, req.GetChatId()); err != nil {
		log.Err(errs.WrapErr(err)).Msg("cleanup chat")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}
		return nil, status.Error(codes.Internal, "failed to cleanup chat")
	}

	return &emptypb.Empty{}, status.Error(codes.OK, "cleaned up chat successfully")
}
