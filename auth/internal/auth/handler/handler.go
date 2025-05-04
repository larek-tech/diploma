package handler

import (
	"context"

	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"go.opentelemetry.io/otel/trace"
)

type authController interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error)
}

// Handler implements authorization on transport level.
type Handler struct {
	pb.UnimplementedAuthServiceServer
	tracer trace.Tracer
	ac     authController
}

// New creates new Handler.
func New(tracer trace.Tracer, ac authController) *Handler {
	return &Handler{
		tracer: tracer,
		ac:     ac,
	}
}
