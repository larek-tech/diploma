package handler

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
)

type roleController interface {
	CreateRole(ctx context.Context, req *pb.CreateRoleRequest, meta *authpb.UserAuthMetadata) (*pb.Role, error)
	GetRole(ctx context.Context, req *pb.GetRoleRequest, meta *authpb.UserAuthMetadata) (*pb.Role, error)
	DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest, meta *authpb.UserAuthMetadata) error
	ListRoles(ctx context.Context, req *pb.ListRolesRequest, meta *authpb.UserAuthMetadata) (*pb.ListRolesResponse, error)
	SetRole(ctx context.Context, req *pb.UpdateRoleRequest, meta *authpb.UserAuthMetadata) error
	RemoveRole(ctx context.Context, req *pb.UpdateRoleRequest, meta *authpb.UserAuthMetadata) error
}

// Handler implements role methods on transport layer.
type Handler struct {
	pb.UnimplementedRoleServiceServer
	rc roleController
}

// New creates new Handler.
func New(rc roleController) *Handler {
	return &Handler{
		rc: rc,
	}
}
