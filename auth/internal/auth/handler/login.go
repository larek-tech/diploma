package handler

import (
	"context"

	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	rescodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login authorizes user with credentials.
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, span := h.tracer.Start(
		ctx,
		"Handler.Login",
		trace.WithAttributes(attribute.String("email", req.GetEmail())),
	)
	defer span.End()

	resp, err := h.ac.Login(ctx, req)
	if err != nil {
		logError(errs.WrapErr(err), span, "login")
		return nil, status.Error(rescodes.Unauthenticated, "failed to login")
	}

	return resp, status.Error(rescodes.OK, "login successful")
}
