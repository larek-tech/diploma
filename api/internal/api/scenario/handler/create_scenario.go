package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// CreateScenario godoc
//
//	@Summary		Create new scenario.
//	@Description	Creates new scenario (group of sources used for RAG vector search).
//	@Tags			scenario
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.CreateScenarioRequest	true	"Input data for creating scenario"
//	@Success		201	{object}	pb.Scenario					"Scenario successfully created"
//	@Failure		400	{object}	string						"Failed to create scenario"
//	@Router			/api/v1/scenario/ [post]
func (h *Handler) CreateScenario(c *fiber.Ctx) error {
	var req pb.CreateScenarioRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.scenarioService.CreateScenario(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrCreateScenario, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
