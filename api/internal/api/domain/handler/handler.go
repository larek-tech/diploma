package handler

import (
	"context"
	"strconv"
	"strings"

	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
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

func pushUserMeta(ctx context.Context, meta *authpb.UserAuthMetadata) context.Context {
	ctx = grpcclient.PushOutMeta(ctx, shared.UserIDHeader, strconv.FormatInt(meta.GetUserId(), 10))
	rolesRaw := meta.GetRoles()
	roles := make([]string, len(rolesRaw))
	for idx := range roles {
		roles[idx] = strconv.FormatInt(rolesRaw[idx], 10)
	}
	ctx = grpcclient.PushOutMeta(ctx, shared.UserRolesHeader, strings.Join(roles, ","))
	return ctx
}
