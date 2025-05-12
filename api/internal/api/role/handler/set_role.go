package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetRole godoc
//
//	@Summary		Set role.
//	@Description	Add new role for user, only for admins.
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.UpdateRoleRequest	true	"Set role for user"
//	@Success		201	{object}	string					"Role successfully set"
//	@Failure		400	{object}	string					"Failed to set role"
//	@Failure		403	{object}	string					"Required admin role"
//	@Router			/api/v1/role/set [put]
func (h *Handler) SetRole(c *fiber.Ctx) error {
	var req pb.UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	_, err := h.roleService.SetRole(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrUpdateRoleForUser, err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}
