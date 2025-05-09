package handler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/chat/internal/auth"
	"github.com/larek-tech/diploma/chat/internal/chat/controller"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RenameChat renames chat title.
func (h *Handler) RenameChat(ctx context.Context, req *pb.RenameChatRequest) (*pb.Chat, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.cc.RenameChat(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("rename chat")
		if errors.Is(err, controller.ErrNoAccessToChat) {
			return nil, status.Error(codes.PermissionDenied, "user doesn't have enough rights")
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "chat not found")
		}
		return nil, status.Error(codes.Internal, "failed to rename chat")
	}

	return resp, status.Error(codes.OK, "renamed chat successfully")
}
