package controller

import (
	"context"
	"time"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdateScenario updates scenario data.
func (ctrl *Controller) UpdateScenario(ctx context.Context, req *pb.UpdateScenarioRequest, meta *authpb.UserAuthMetadata) (*pb.Scenario, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.UpdateScenario",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.Int64("scenarioID", req.GetScenarioId()),
			attribute.String("llm model", req.GetModelName()),
		),
	)
	defer span.End()

	scenario, err := ctrl.sr.GetScenarioByID(ctx, req.GetScenarioId(), meta.GetUserId())
	if err != nil {
		return nil, errs.WrapErr(err)
	}

	scenario.UseMultiquery = req.GetUseMultiquery()
	scenario.NQueries = req.GetNQueries()
	scenario.QueryModelName = req.GetQueryModelName()
	scenario.UseRerank = req.GetUseRerank()
	scenario.RerankerModelName = req.GetRerankerModel()
	scenario.RerankerMaxLength = req.GetRerankerMaxLength()
	scenario.RerankerTopK = req.GetRerankerTopK()
	scenario.LlmModelName = req.GetModelName()
	scenario.Temperature = req.GetTemperature()
	scenario.TopK = req.GetModelTopK()
	scenario.TopP = req.GetTopP()
	scenario.SystemPrompt = req.GetSystemPrompt()
	scenario.TopN = req.GetTopN()
	scenario.Threshold = req.GetThreshold()
	scenario.SearchByQuery = req.GetSearchByQuery()
	scenario.UpdatedAt = time.Now()

	if err = ctrl.sr.UpdateScenario(ctx, scenario, meta.GetUserId()); err != nil {
		return nil, errs.WrapErr(err)
	}

	return scenario.ToProto(), nil
}
