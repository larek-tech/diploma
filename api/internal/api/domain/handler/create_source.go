package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// CreateSource godoc
// @Summary Create new source.
// @Description Creates new domain source used for RAG vector search.
// @Tags domain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param req body pb.CreateSourceRequest true "Input data for creating source"
// @Success 201 {object} pb.CreateSourceResponse "Source successfully created"
// @Failure 400 {object} string "Failed to create source"
// @Router /api/v1/domain/ [post]
func (h *Handler) CreateSource(c *fiber.Ctx) error {
	var req pb.CreateSourceRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	userID, ok := c.Locals(shared.UserIDKey).(int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user ID in context")
	}
	req.UserId = userID

	resp, err := h.domainService.CreateSource(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrCreateSource, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
