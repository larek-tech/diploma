package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateRole godoc
//
//	@Summary		Create new role.
//	@Description	Create new role, only for admins.
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.CreateRoleRequest	true	"Input data for creating role"
//	@Success		201	{object}	pb.Role					"Role successfully created"
//	@Failure		400	{object}	string					"Failed to create role"
//	@Failure		403	{object}	string					"Required admin role"
//	@Router			/api/v1/role/ [post]
func (h *Handler) CreateRole(c *fiber.Ctx) error {
	var req pb.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.roleService.CreateRole(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(shared.ErrCreateRole, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
