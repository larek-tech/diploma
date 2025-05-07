package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPermittedRoles godoc
//
//	@Summary		Get permitted roles.
//	@Description	Returns list of roles permitted to source.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			sourceID	path		int						true	"Requested source ID"
//	@Success		200			{object}	pb.GetSourceResponse	"Permitted roles"
//	@Failure		404			{object}	string					"Source not found"
//	@Router			/api/v1/source/permissions/roles/{id} [get]
func (h *Handler) GetPermittedRoles(c *fiber.Ctx) error {
	var req pb.GetResourcePermissionsRequest

	sourceID, err := c.ParamsInt(sourceIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ResourceId = int64(sourceID)

	resp, err := h.sourceService.GetPermittedRoles(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
