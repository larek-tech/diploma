package handler

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

const (
	sourceIDParam = "id"
	offsetParam   = "offset"
	limitParam    = "limit"
)

// Handler implements source methods on transport level.
type Handler struct {
	sourceService pb.SourceServiceClient
	tracer        trace.Tracer
}

// New creates new Handler.
func New(sourceService pb.SourceServiceClient, tracer trace.Tracer) *Handler {
	return &Handler{
		sourceService: sourceService,
		tracer:        tracer,
	}
}
