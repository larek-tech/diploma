package handler

import (
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	domainpb "github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

const (
	chatIDParam  = "id"
	queryIDParam = "id"
	offsetParam  = "offset"
	limitParam   = "limit"
)

// Handler implements chat methods on transport level.
type Handler struct {
	chatService     pb.ChatServiceClient
	authService     authpb.AuthServiceClient
	mlService       domainpb.MLServiceClient
	scenarioService domainpb.ScenarioServiceClient
	domainService   domainpb.DomainServiceClient
	sourceService   domainpb.SourceServiceClient
	tracer          trace.Tracer
}

// New creates new Handler.
func New(
	chatService pb.ChatServiceClient,
	authService authpb.AuthServiceClient,
	mlService domainpb.MLServiceClient,
	scenarioService domainpb.ScenarioServiceClient,
	domainService domainpb.DomainServiceClient,
	sourceService domainpb.SourceServiceClient,
	tracer trace.Tracer,
) *Handler {
	return &Handler{
		chatService:     chatService,
		authService:     authService,
		mlService:       mlService,
		scenarioService: scenarioService,
		domainService:   domainService,
		sourceService:   sourceService,
		tracer:          tracer,
	}
}
