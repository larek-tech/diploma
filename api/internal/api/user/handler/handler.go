package handler

import "github.com/larek-tech/diploma/api/internal/domain/pb"

const (
	userIDParam = "id"
	offsetParam   = "offset"
	limitParam    = "limit"
)

// Handler implements user methods on transport layer.
type Handler struct {
	userService pb.UserServiceClient
}

// New creates new Handler.
func New(userService pb.UserServiceClient) *Handler {
	return &Handler{
		userService: userService,
	}
}
