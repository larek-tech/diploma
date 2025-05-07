package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/api/domain"
	dh "github.com/larek-tech/diploma/api/internal/api/domain/handler"
	"github.com/larek-tech/diploma/api/internal/api/scenario"
	sch "github.com/larek-tech/diploma/api/internal/api/scenario/handler"
	"github.com/larek-tech/diploma/api/internal/api/source"
	sh "github.com/larek-tech/diploma/api/internal/api/source/handler"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// SetupRoutes maps api routes.
func SetupRoutes(api fiber.Router, tracer trace.Tracer, domainConn *grpc.ClientConn) {
	sourceRouter := api.Group("/source")
	sourceHandler := sh.New(pb.NewSourceServiceClient(domainConn), tracer)
	source.SetupRoutes(sourceRouter, sourceHandler)

	domainRouter := api.Group("/domain")
	domainHandler := dh.New(pb.NewDomainServiceClient(domainConn), tracer)
	domain.SetupRoutes(domainRouter, domainHandler)

	scenarioRouter := api.Group("/scenario")
	scenarioHandler := sch.New(pb.NewScenarioServiceClient(domainConn), tracer)
	scenario.SetupRoutes(scenarioRouter, scenarioHandler)
}
