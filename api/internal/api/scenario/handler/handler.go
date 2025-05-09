package handler

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
)

const (
	scenarioIDParam = "id"
	offsetParam     = "offset"
	limitParam      = "limit"
)

// Handler implements scenario methods on transport level.
type Handler struct {
	scenarioService pb.ScenarioServiceClient
}

// New creates new Handler.
func New(scenarioService pb.ScenarioServiceClient) *Handler {
	return &Handler{
		scenarioService: scenarioService,
	}
}
