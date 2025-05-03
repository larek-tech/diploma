package handler

import (
	"context"
	"strconv"

	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	grpcclient "github.com/yogenyslav/pkg/grpc_client"
	"go.opentelemetry.io/otel/trace"
)

const (
	sourceIDParam = "id"
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

func pushUserID(ctx context.Context, userID int64) context.Context {
	return grpcclient.PushOutMeta(ctx, shared.UserIDHeader, strconv.FormatInt(userID, 10))
}
