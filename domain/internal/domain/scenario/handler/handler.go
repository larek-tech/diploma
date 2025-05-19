package handler

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

type scenarioController interface {
	CreateScenario(ctx context.Context, req *pb.CreateScenarioRequest, meta *authpb.UserAuthMetadata) (*pb.Scenario, error)
	GetScenario(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) (*pb.Scenario, error)
	GetDefaultScenario(ctx context.Context, req *pb.GetDefaultScenarioRequest, meta *authpb.UserAuthMetadata) (*pb.Scenario, error)
	UpdateScenario(ctx context.Context, req *pb.UpdateScenarioRequest, meta *authpb.UserAuthMetadata) (*pb.Scenario, error)
	DeleteScenario(ctx context.Context, sourceID int64, meta *authpb.UserAuthMetadata) error
	ListScenarios(ctx context.Context, req *pb.ListScenariosRequest, meta *authpb.UserAuthMetadata) (*pb.ListScenariosResponse, error)
	ListScenariosByDomain(ctx context.Context, req *pb.ListScenariosByDomainRequest, meta *authpb.UserAuthMetadata) (*pb.ListScenariosResponse, error)
}

// Handler implements source methods on transport level.
type Handler struct {
	pb.UnimplementedScenarioServiceServer
	sc     scenarioController
	tracer trace.Tracer
}

// New creates new Handler.
func New(sc scenarioController, tracer trace.Tracer) *Handler {
	return &Handler{
		sc:     sc,
		tracer: tracer,
	}
}
