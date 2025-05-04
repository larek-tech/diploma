package handler

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

type domainController interface {
	CreateDomain(ctx context.Context, req *pb.CreateDomainRequest, meta *authpb.UserAuthMetadata) (*pb.Domain, error)
	GetDomain(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) (*pb.GetDomainResponse, error)
	UpdateDomain(ctx context.Context, req *pb.UpdateDomainRequest, meta *authpb.UserAuthMetadata) (*pb.Domain, error)
	DeleteDomain(ctx context.Context, domainID int64, meta *authpb.UserAuthMetadata) error
	ListDomains(ctx context.Context, req *pb.ListDomainsRequest, meta *authpb.UserAuthMetadata) (*pb.ListDomainsResponse, error)
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
