package handler

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

const (
	domainIDParam = "id"
	offsetParam   = "offset"
	limitParam    = "limit"
)

// Handler implements domain methods on transport level.
type Handler struct {
	domainService pb.DomainServiceClient
	tracer        trace.Tracer
}

// New creates new Handler.
func New(domainService pb.DomainServiceClient, tracer trace.Tracer) *Handler {
	return &Handler{
		domainService: domainService,
		tracer:        tracer,
	}
}
