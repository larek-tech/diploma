package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetDomain godoc
//
//	@Summary		Get domain.
//	@Description	Returns information about domain.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int			true	"Requested domain ID"
//	@Success		200	{object}	pb.Domain	"Domain"
//	@Failure		400	{object}	string		"Failed to get domain"
//	@Failure		404	{object}	string		"Domain not found"
//	@Router			/api/v1/domain/{id} [get]
func (h *Handler) GetDomain(c *fiber.Ctx) error {
	var req pb.GetDomainRequest

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.DomainId = int64(domainID)

	resp, err := h.domainService.GetDomain(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrDomainNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrGetDomain, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
