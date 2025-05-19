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

// ListScenariosByDomain returns paginated list of scenarios by specified domain.
func (h *Handler) ListScenariosByDomain(ctx context.Context, req *pb.ListScenariosByDomainRequest) (*pb.ListScenariosResponse, error) {
	meta, err := auth.GetUserMeta(ctx)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("get user meta")
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	resp, err := h.sc.ListScenariosByDomain(ctx, req, meta)
	if err != nil {
		log.Err(errs.WrapErr(err)).Msg("list scenarios by domain")
		return nil, status.Error(codes.Internal, "failed to list scenarios by domain")
	}

	return resp, status.Error(codes.OK, "got scenarios list by domain successfully")
}
