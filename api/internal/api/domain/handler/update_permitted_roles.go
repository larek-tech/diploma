package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePermittedRoles godoc
//
//	@Summary		Update permitted roles.
//	@Description	Updates list of roles permitted to domain.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int					true	"Requested domain ID"
//	@Param			req	body		pb.PermittedRoles	true	"New list of permitted roles"
//	@Success		200	{object}	pb.PermittedRoles	"Updated roles permissions"
//	@Failure		404	{object}	string				"Domain not found"
//	@Router			/api/v1/domain/permissions/roles/{id} [put]
func (h *Handler) UpdatePermittedRoles(c *fiber.Ctx) error {
	var req pb.PermittedRoles
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ResourceId = int64(domainID)

	resp, err := h.domainService.UpdatePermittedRoles(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
