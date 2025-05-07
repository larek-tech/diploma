package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateScenario creates new scenario.
func (h *Handler) CreateScenario(ctx context.Context, req *pb.CreateScenarioRequest) (*pb.Scenario, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.CreateScenario(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create scenario")
		return nil, status.Error(codes.Internal, "failed to create scenario")
	}

	return resp, status.Error(codes.OK, "scenario created successfully")
}
