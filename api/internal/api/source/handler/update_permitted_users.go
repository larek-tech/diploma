package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdatePermittedUsers godoc
//
//	@Summary		Update permitted users.
//	@Description	Updates list of users permitted to source.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			sourceID	path		int					true	"Requested source ID"
//	@Param			req			body		pb.PermittedUsers	true	"New list of permitted users"
//	@Success		200			{object}	pb.PermittedUsers	"Updated users permissions"
//	@Failure		404			{object}	string				"Source not found"
//	@Router			/api/v1/source/permissions/users/{id} [put]
func (h *Handler) UpdatePermittedUsers(c *fiber.Ctx) error {
	var req pb.PermittedUsers
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	sourceID, err := c.ParamsInt(sourceIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ResourceId = int64(sourceID)

	resp, err := h.sourceService.UpdatePermittedUsers(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
