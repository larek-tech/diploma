package controller

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// GetScenario returns scenario by id.
func (ctrl *Controller) GetScenario(ctx context.Context, scenarioID int64, meta *authpb.UserAuthMetadata) (*pb.Scenario, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.GetScenario",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("scenarioID", scenarioID),
		),
	)
	defer span.End()

	scenario, err := ctrl.sr.GetScenarioByID(ctx, scenarioID, meta.GetUserId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	return scenario.ToProto(), nil
}
