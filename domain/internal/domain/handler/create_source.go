package handler

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateSource creates new source.
func (h *Handler) CreateSource(ctx context.Context, req *pb.Source) (*pb.Source, error) {
	meta, err := getUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	if updateParams := req.GetUpdateParams(); updateParams != nil {
		switch {
		case updateParams.EveryPeriod != nil && *updateParams.EveryPeriod <= 0:
			return nil, status.Error(codes.InvalidArgument, "update period must be positive")
		case updateParams.GetCron() != nil:
			cron := updateParams.GetCron()
			if cron.GetDayOfWeek() < 0 || cron.Month < 0 || cron.DayOfMonth < 0 || cron.Hour < 0 || cron.Minute < 0 {
				return nil, status.Error(codes.InvalidArgument, "cron update params must be non-negative")
			}
		}
	}

	resp, err := h.sc.CreateSource(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("create source")
		return nil, status.Error(codes.Internal, "failed to create source")
	}

	return resp, status.Error(codes.OK, "source created successfully")
}
