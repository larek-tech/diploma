package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListScenarios returns the paginated list of available scenarios.
func (ctrl *Controller) ListScenarios(ctx context.Context, req *pb.ListScenariosRequest, meta *authpb.UserAuthMetadata) (*pb.ListScenariosResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.ListScenarios",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("offest", int64(req.GetOffset())),
			attribute.Int64("limit", int64(req.GetLimit())),
		),
	)
	defer span.End()

	scenariosDao, err := ctrl.sr.ListScenarios(ctx, meta.GetUserId(), req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	scenarios := make([]*pb.Scenario, len(scenariosDao))
	for idx := range scenariosDao {
		scenarios[idx] = scenariosDao[idx].ToProto()
	}

	return &pb.ListScenariosResponse{Scenarios: scenarios}, nil
}
