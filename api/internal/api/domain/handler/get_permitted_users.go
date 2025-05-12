package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPermittedUsers godoc
//
//	@Summary		Get permitted users.
//	@Description	Returns list of users permitted to domain.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int			true	"Requested domain ID"
//	@Success		200	{object}	pb.Domain	"Permitted users"
//	@Failure		404	{object}	string		"Domain not found"
//	@Router			/api/v1/domain/permissions/users/{id} [get]
func (h *Handler) GetPermittedUsers(c *fiber.Ctx) error {
	var req pb.GetResourcePermissionsRequest

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ResourceId = int64(domainID)

	resp, err := h.domainService.GetPermittedUsers(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
