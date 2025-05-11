package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetDefaultScenario returns default scenario for domain.
func (ctrl *Controller) GetDefaultScenario(
	ctx context.Context,
	req *pb.GetDefaultScenarioRequest,
	meta *authpb.UserAuthMetadata,
) (*pb.Scenario, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetDefaultScenario",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("title", req.GetDefaultTitle()),
		),
	)
	defer span.End()

	scenario, err := ctrl.sr.GetDefaultScenario(ctx, req.GetDefaultTitle(), meta.GetUserId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return scenario.ToProto(), nil
}
