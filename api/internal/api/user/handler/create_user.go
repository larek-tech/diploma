package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser godoc
//
//	@Summary		Create new user.
//	@Description	Create new user, only for admins.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.CreateUserRequest	true	"Input data for creating user"
//	@Success		201	{object}	pb.User					"User successfully created"
//	@Failure		400	{object}	string					"Failed to create user"
//	@Failure		403	{object}	string					"Required admin role"
//	@Router			/api/v1/user/ [post]
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req pb.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.userService.CreateUser(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrCreateUser, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
