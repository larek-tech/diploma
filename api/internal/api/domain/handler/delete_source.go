package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// DeleteSource godoc
// @Summary Delete source.
// @Description Delete domain source by ID.
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sourceID path int true "Source ID"
// @Success 204 {object} string "Source deleted"
// @Failure 400 {object} string "Failed to delete source"
// @Router /api/v1/domain/{id} [delete]
func (h *Handler) DeleteSource(c *fiber.Ctx) error {
	var req pb.DeleteSourceRequest

	userID, ok := c.Locals(shared.UserIDKey).(int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user ID in context")
	}

	sourceID, err := c.ParamsInt(sourceIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidPath, err.Error())
	}

	req.UserId = userID
	req.SourceId = int64(sourceID)

	_, err = h.domainService.DeleteSource(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrDeleteSource, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
