package handler

import (
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

type domainController interface {
}

// Handler implements domain methods on transport level.
type Handler struct {
	pb.UnimplementedDomainServiceServer
	dc     domainController
	tracer trace.Tracer
}

// New creates new Handler.
func New(dc domainController, tracer trace.Tracer) *Handler {
	return &Handler{
		tracer: tracer,
		dc:     dc,
	}
}
