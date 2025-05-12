package handler

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

type domainController interface {
	CreateDomain(ctx context.Context, req *pb.CreateDomainRequest, meta *authpb.UserAuthMetadata) (*pb.Domain, error)
	GetDomain(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) (*pb.Domain, error)
	UpdateDomain(ctx context.Context, req *pb.UpdateDomainRequest, meta *authpb.UserAuthMetadata) (*pb.Domain, error)
	DeleteDomain(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) error
	ListDomains(ctx context.Context, req *pb.ListDomainsRequest, meta *authpb.UserAuthMetadata) (*pb.ListDomainsResponse, error)
	GetPermittedRoles(ctx context.Context, req *pb.GetResourcePermissionsRequest, meta *authpb.UserAuthMetadata) (*pb.PermittedRoles, error)
	GetPermittedUsers(ctx context.Context, req *pb.GetResourcePermissionsRequest, meta *authpb.UserAuthMetadata) (*pb.PermittedUsers, error)
	UpdatePermittedRoles(ctx context.Context, req *pb.PermittedRoles, meta *authpb.UserAuthMetadata) (*pb.PermittedRoles, error)
	UpdatePermittedUsers(ctx context.Context, req *pb.PermittedUsers, meta *authpb.UserAuthMetadata) (*pb.PermittedUsers, error)
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
		dc:     dc,
		tracer: tracer,
	}
}
