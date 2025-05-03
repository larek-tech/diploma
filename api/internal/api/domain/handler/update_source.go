package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// UpdateSource godoc
//
//	@Summary		Update source.
//	@Description	Update source information or params for update job.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			sourceID	path		int						true	"Source ID"
//	@Param			req			body		pb.UpdateSourceRequest	true	"Update params"
//	@Success		200			{object}	pb.Source				"Source updated"
//	@Failure		400			{object}	string					"Failed to update source"
//	@Router			/api/v1/domain/{id} [put]
func (h *Handler) UpdateSource(c *fiber.Ctx) error {
	var req pb.UpdateSourceRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	userID, ok := c.Locals(shared.UserIDKey).(int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user ID in context")
	}

	sourceID, err := c.ParamsInt(sourceIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}

	ctx := pushUserID(c.UserContext(), userID)
	req.SourceId = int64(sourceID)

	resp, err := h.domainService.UpdateSource(ctx, &req)
	if err != nil {
		return errs.WrapErr(shared.ErrUpdateSource, err.Error())
	}

	return c.Status(fiber.StatusNoContent).JSON(resp)
}
