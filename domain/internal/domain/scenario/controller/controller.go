package controller

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/scenario/model"
	"go.opentelemetry.io/otel/trace"
)

type scenarioRepo interface {
	InsertScenario(ctx context.Context, s model.ScenarioDao) (int64, error)
	GetScenarioByID(ctx context.Context, id, userID int64) (model.ScenarioDao, error)
	GetDefaultScenario(ctx context.Context, title string, userID int64) (model.ScenarioDao, error)
	UpdateScenario(ctx context.Context, s model.ScenarioDao, userID int64) error
	DeleteScenario(ctx context.Context, id, userID int64) error
	ListScenarios(ctx context.Context, userID int64, offset, limit uint64) ([]model.ScenarioDao, error)
}

// Controller implements scenario methods on logic layer.
type Controller struct {
	sr     scenarioRepo
	tracer trace.Tracer
}

// New creates new Controller.
func New(sr scenarioRepo, tracer trace.Tracer) *Controller {
	return &Controller{
		sr:     sr,
		tracer: tracer,
	}
}
