package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetScenario godoc
//
//	@Summary		Get scenario.
//	@Description	Returns information about scenario.
//	@Tags			scenario
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			scenarioID	path		int						true	"Requested scenario ID"
//	@Success		200			{object}	pb.GetScenarioResponse	"Scenario"
//	@Failure		400			{object}	string					"Failed to get scenario"
//	@Failure		404			{object}	string					"Scenario not found"
//	@Router			/api/v1/scenario/{id} [get]
func (h *Handler) GetScenario(c *fiber.Ctx) error {
	var req pb.GetScenarioRequest

	scenarioID, err := c.ParamsInt(scenarioIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ScenarioId = int64(scenarioID)

	resp, err := h.scenarioService.GetScenario(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrScenarioNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrGetScenario, err.Error())
	}

	if resp.GetScenario() == nil {
		return errs.WrapErr(shared.ErrScenarioNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
