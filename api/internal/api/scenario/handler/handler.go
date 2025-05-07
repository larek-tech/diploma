package handler

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

const (
	scenarioIDParam = "id"
	offsetParam     = "offset"
	limitParam      = "limit"
)

// Handler implements scenario methods on transport level.
type Handler struct {
	scenarioService pb.ScenarioServiceClient
	tracer          trace.Tracer
}

// New creates new Handler.
func New(scenarioService pb.ScenarioServiceClient, tracer trace.Tracer) *Handler {
	return &Handler{
		scenarioService: scenarioService,
		tracer:          tracer,
	}
}
