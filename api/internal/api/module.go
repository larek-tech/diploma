package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/api/domain"
	"github.com/larek-tech/diploma/api/internal/api/domain/handler"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// SetupRoutes maps api routes.
func SetupRoutes(api fiber.Router, tracer trace.Tracer, domainConn *grpc.ClientConn) {
	domainHandler := handler.New(pb.NewDomainServiceClient(domainConn), tracer)
	domainRouter := api.Group("/domain")
	domain.SetupRoutes(domainRouter, domainHandler)
}
