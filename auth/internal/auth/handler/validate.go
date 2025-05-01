package handler

import (
	"context"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	rescodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Validate validates provided access token.
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	ctx, span := h.tracer.Start(ctx, "Handler.Validate")
	defer span.End()

	resp, err := h.ac.Validate(ctx, req)
	if err != nil {
		logError(errs.WrapErr(err), span)
		return nil, status.Error(rescodes.Unauthenticated, "invalid credentials")
	}

	span.SetAttributes(
		attribute.Int64("userID", resp.GetMeta().UserId),
		attribute.Int64Slice("roleIDs", resp.GetMeta().Roles),
	)

	return resp, status.Error(rescodes.OK, "validation successful")
}
