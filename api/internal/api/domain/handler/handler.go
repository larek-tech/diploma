package handler

import (
	"fmt"

	"github.com/larek-tech/diploma/api/internal/domain/pb"
)

const (
	domainIDParam = "id"
	offsetParam   = "offset"
	limitParam    = "limit"

	defaultTitlePattern = "%s-%d (сценарий по умолчанию)"
	optimalTitlePattern = "%s-%d (оптимальные параметры)"
)

// Handler implements domain methods on transport level.
type Handler struct {
	domainService   pb.DomainServiceClient
	scenarioService pb.ScenarioServiceClient
	sourceService   pb.SourceServiceClient
	mlService       pb.MLServiceClient
}

// New creates new Handler.
func New(
	domainService pb.DomainServiceClient,
	scenarioService pb.ScenarioServiceClient,
	sourceService pb.SourceServiceClient,
	mlService pb.MLServiceClient,
) *Handler {
	return &Handler{
		domainService:   domainService,
		scenarioService: scenarioService,
		sourceService:   sourceService,
		mlService:       mlService,
	}
}

func domainDefaultTitle(domainTitle string, domainID int64) string {
	return fmt.Sprintf(defaultTitlePattern, domainTitle, domainID)
}

func domainOptimalTitle(domainTitle string, domainID int64) string {
	return fmt.Sprintf(optimalTitlePattern, domainTitle, domainID)
}
