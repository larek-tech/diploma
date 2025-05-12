package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUser godoc
//
//	@Summary		Delete user.
//	@Description	Delete user by ID, only for admins.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int		true	"User ID"
//	@Success		204	{object}	string	"User deleted"
//	@Failure		400	{object}	string	"Failed to delete user"
//	@Failure		403	{object}	string	"Required admin role"
//	@Router			/api/v1/user/{id} [delete]
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	var req pb.DeleteUserRequest

	userID, err := c.ParamsInt(userIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.UserId = int64(userID)

	_, err = h.userService.DeleteUser(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrDeleteUser, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
