package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateDomain godoc
//
//	@Summary		Update domain.
//	@Description	Update domain information.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int						true	"Domain ID"
//	@Param			req	body		pb.UpdateDomainRequest	true	"Update params"
//	@Success		200	{object}	pb.Domain				"Domain updated"
//	@Failure		400	{object}	string					"Failed to update domain"
//	@Failure		404	{object}	string					"Domain not found"
//	@Router			/api/v1/domain/{id} [put]
func (h *Handler) UpdateDomain(c *fiber.Ctx) error {
	var req pb.UpdateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.DomainId = int64(domainID)

	domain, err := h.domainService.UpdateDomain(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrDomainNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrUpdateDomain, err.Error())
	}

	userID := c.Locals(shared.UserIDKey).(int64)
	roles := c.Locals(shared.UserRolesKey).([]int64)

	resp, err := h.checkDefaultScenario(c.UserContext(), domain, userID, roles)
	if err != nil {
		return errs.WrapErr(shared.ErrUpdateDomain, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
