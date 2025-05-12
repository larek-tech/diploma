package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteScenario godoc
//
//	@Summary		Delete scenario.
//	@Description	Delete scenario by ID.
//	@Tags			scenario
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int		true	"Scenario ID"
//	@Success		204	{object}	string	"Scenario deleted"
//	@Failure		400	{object}	string	"Failed to delete scenario"
//	@Failure		404	{object}	string	"Scenario not found"
//	@Router			/api/v1/scenario/{id} [delete]
func (h *Handler) DeleteScenario(c *fiber.Ctx) error {
	var req pb.DeleteScenarioRequest

	scenarioID, err := c.ParamsInt(scenarioIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ScenarioId = int64(scenarioID)

	_, err = h.scenarioService.DeleteScenario(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrScenarioNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrDeleteScenario, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
