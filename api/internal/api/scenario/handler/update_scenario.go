package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateScenario godoc
//
//	@Summary		Update scenario.
//	@Description	Update scenario information.
//	@Tags			scenario
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			scenarioID	path		int							true	"Scenario ID"
//	@Param			req			body		pb.UpdateScenarioRequest	true	"Update params"
//	@Success		200			{object}	pb.Scenario					"Scenario updated"
//	@Failure		400			{object}	string						"Failed to update scenario"
//	@Failure		404			{object}	string						"Scenario not found"
//	@Router			/api/v1/scenario/{id} [put]
func (h *Handler) UpdateScenario(c *fiber.Ctx) error {
	var req pb.UpdateScenarioRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	scenarioID, err := c.ParamsInt(scenarioIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ScenarioId = int64(scenarioID)

	resp, err := h.scenarioService.UpdateScenario(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrScenarioNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrUpdateScenario, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
