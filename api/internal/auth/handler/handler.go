package handler

import (
	"github.com/larek-tech/diploma/api/internal/auth/pb"
	"go.opentelemetry.io/otel/trace"
)

// Handler implements authorization on transport level.
type Handler struct {
	authService pb.AuthServiceClient
	tracer      trace.Tracer
}

// New creates new Handler.
func New(authService pb.AuthServiceClient, tracer trace.Tracer) *Handler {
	return &Handler{
		authService: authService,
		tracer:      tracer,
	}
}
