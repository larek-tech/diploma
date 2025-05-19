package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ListScenariosByDomain returns paginated list of scenarios by specified domain.
func (ctrl *Controller) ListScenariosByDomain(
	ctx context.Context,
	req *pb.ListScenariosByDomainRequest,
	meta *authpb.UserAuthMetadata,
) (*pb.ListScenariosResponse, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"ListScenariosByDomain",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("domainID", req.GetDomainId()),
			attribute.Int64("offset", int64(req.GetOffset())),
			attribute.Int64("limit", int64(req.GetLimit())),
		),
	)
	defer span.End()

	scenariosDB, err := ctrl.sr.ListScenariosByDomain(
		ctx,
		req.GetDomainId(),
		req.GetOffset(),
		req.GetLimit(),
	)
	if err != nil {
		return nil, errs.WrapErr(err, "list scenarios by domain")
	}

	scenarios := make([]*pb.Scenario, len(scenariosDB))
	for idx := range scenariosDB {
		scenarios[idx] = scenariosDB[idx].ToProto()
	}

	return &pb.ListScenariosResponse{
		Scenarios: scenarios,
	}, nil
}
