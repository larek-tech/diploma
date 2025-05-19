package controller

import (
	"context"
	"time"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// CreateScenario creates new scenario record.
func (ctrl *Controller) CreateScenario(ctx context.Context, req *pb.CreateScenarioRequest, meta *authpb.UserAuthMetadata) (*pb.Scenario, error) {
	ctx, span := ctrl.tracer.Start(
		ctx,
		"Controller.CreateScenario",
		trace.WithAttributes(
			attribute.Int64("userID", meta.GetUserId()),
			attribute.String("llm model", req.GetModel().GetModelName()),
		),
	)
	defer span.End()

	multiQuery := req.GetMultiQuery()
	reranker := req.GetReranker()
	llm := req.GetModel()
	vectorSearch := req.GetVectorSearch()

	scenario := model.ScenarioDao{
		Title:             req.GetTitle(),
		UserID:            meta.GetUserId(),
		DomainID:          req.GetDomainId(),
		ContextSize:       req.GetContextSize(),
		UseMultiquery:     multiQuery.GetUseMultiquery(),
		NQueries:          multiQuery.GetNQueries(),
		QueryModelName:    multiQuery.GetQueryModelName(),
		UseRerank:         reranker.GetUseRerank(),
		RerankerModelName: reranker.GetRerankerModel(),
		RerankerMaxLength: reranker.GetRerankerMaxLength(),
		RerankerTopK:      reranker.GetTopK(),
		LlmModelName:      llm.GetModelName(),
		Temperature:       llm.GetTemperature(),
		TopK:              llm.GetTopK(),
		TopP:              llm.GetTopP(),
		SystemPrompt:      llm.GetSystemPrompt(),
		TopN:              vectorSearch.GetTopN(),
		Threshold:         vectorSearch.GetThreshold(),
		SearchByQuery:     vectorSearch.GetSearchByQuery(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	scenarioID, err := ctrl.sr.InsertScenario(ctx, scenario)
	if err != nil {
		return nil, errs.WrapErr(err)
	}
	scenario.ID = scenarioID

	return scenario.ToProto(), nil
}
