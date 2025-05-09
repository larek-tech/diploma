package handler

import (
	"context"

	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	rescodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login authorizes user with credentials.
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := h.ac.Login(ctx, req)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("failed to login")
		return nil, status.Error(rescodes.Unauthenticated, "failed to login")
	}

	return resp, status.Error(rescodes.OK, "login successful")
}
