package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteRole godoc
//
//	@Summary		Delete role.
//	@Description	Delete role by ID, only for admins.
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int		true	"Role ID"
//	@Success		204	{object}	string	"Role deleted"
//	@Failure		400	{object}	string	"Failed to delete role"
//	@Failure		403	{object}	string	"Required admin role"
//	@Router			/api/v1/role/{id} [delete]
func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	var req pb.DeleteRoleRequest

	roleID, err := c.ParamsInt(roleIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.RoleId = int64(roleID)

	_, err = h.roleService.DeleteRole(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrDeleteRole, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
