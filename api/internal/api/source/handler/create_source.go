package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// CreateSource godoc
//
//	@Summary		Create new source.
//	@Description	Creates new source used for RAG vector search.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.CreateSourceRequest	true	"Input data for creating source"
//	@Success		201	{object}	pb.Source				"Source successfully created"
//	@Failure		400	{object}	string					"Failed to create source"
//	@Router			/api/v1/source/ [post]
func (h *Handler) CreateSource(c *fiber.Ctx) error {
	var req pb.CreateSourceRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.sourceService.CreateSource(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrCreateSource, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
