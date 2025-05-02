package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// ListSources godoc
// @Summary List sources.
// @Description List sources to which user has access.
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} pb.ListSourcesResponse "List of sources"
// @Failure 400 {object} string "Failed to list sources"
// @Router /api/v1/domain/list [get]
func (h *Handler) ListSources(c *fiber.Ctx) error {
	var req pb.ListSourcesRequest

	userID, ok := c.Locals(shared.UserIDKey).(int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user ID in context")
	}
	req.UserId = userID

	resp, err := h.domainService.ListSources(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrListSources, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
