package handler

import (
	"context"

	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	rescodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validate validates provided access token.
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	resp, err := h.ac.Validate(ctx, req)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("validate")
		return nil, status.Error(rescodes.Unauthenticated, "invalid credentials")
	}

	return resp, status.Error(rescodes.OK, "validation successful")
}
