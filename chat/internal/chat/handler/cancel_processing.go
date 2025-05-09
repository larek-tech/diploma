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

// CancelProcessing cancels query processing by id.
func (h *Handler) CancelProcessing(ctx context.Context, req *pb.CancelProcessingRequest) (*emptypb.Empty, error) {
	_, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if err = h.cc.CancelProcessing(ctx, req); err != nil {
		log.Err(errs.WrapErr(err)).Msg("cancel processing")
		return nil, status.Errorf(codes.Internal, "failed canceling processing")
	}

	return &emptypb.Empty{}, nil
}
