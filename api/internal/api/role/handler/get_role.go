package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRole godoc
//
//	@Summary		Get role.
//	@Description	Returns information about role.
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int		true	"Requested role ID"
//	@Success		200	{object}	pb.Role	"Role"
//	@Failure		400	{object}	string	"Failed to get role"
//	@Failure		403	{object}	string	"Required admin role"
//	@Router			/api/v1/role/{id} [get]
func (h *Handler) GetRole(c *fiber.Ctx) error {
	var req pb.GetRoleRequest

	roleID, err := c.ParamsInt(roleIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.RoleId = int64(roleID)

	resp, err := h.roleService.GetRole(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrGetRole, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
