package handler

import (
	"context"

	authpb "github.com/larek-tech/diploma/domain/internal/auth/pb"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
)

type userController interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest, meta *authpb.UserAuthMetadata) (*pb.User, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest, meta *authpb.UserAuthMetadata) (*pb.User, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest, meta *authpb.UserAuthMetadata) error
	ListUsers(ctx context.Context, req *pb.ListUsersRequest, meta *authpb.UserAuthMetadata) (*pb.ListUsersResponse, error)
}

// Handler implements user methods on transport layer.
type Handler struct {
	pb.UnimplementedUserServiceServer
	uc userController
}

// New creates new Handler.
func New(uc userController) *Handler {
	return &Handler{
		uc: uc,
	}
}
