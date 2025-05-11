package handler

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/larek-tech/diploma/domain/internal/auth"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetScenario returns scenario by id.
func (h *Handler) GetScenario(ctx context.Context, req *pb.GetScenarioRequest) (*pb.Scenario, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.GetScenario(ctx, req.GetScenarioId(), meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get scenario")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "scenario not found")
		} else {
			return nil, status.Error(codes.Internal, "failed to get scenario")
		}
	}

	return resp, status.Error(codes.OK, "got scenario successfully")
}
