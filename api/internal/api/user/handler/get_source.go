package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser godoc
//
//	@Summary		Get user.
//	@Description	Returns information about user.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int		true	"Requested user ID"
//	@Success		200	{object}	pb.User	"User"
//	@Failure		400	{object}	string	"Failed to get user"
//	@Failure		403	{object}	string	"Required admin role"
//	@Router			/api/v1/user/{id} [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	var req pb.GetUserRequest

	userID, err := c.ParamsInt(userIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.UserId = int64(userID)

	resp, err := h.userService.GetUser(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrGetUser, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
