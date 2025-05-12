package handler

import "github.com/larek-tech/diploma/api/internal/domain/pb"

const (
	roleIDParam = "id"
	offsetParam = "offset"
	limitParam  = "limit"
)

// Handler implements role methods on transport layer.
type Handler struct {
	roleService pb.RoleServiceClient
}

// New creates new Handler.
func New(roleService pb.RoleServiceClient) *Handler {
	return &Handler{
		roleService: roleService,
	}
}
