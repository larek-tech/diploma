package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteDomain godoc
//
//	@Summary		Delete domain.
//	@Description	Delete domain by ID.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			domainID	path		int		true	"Domain ID"
//	@Success		204			{object}	string	"Domain deleted"
//	@Failure		400			{object}	string	"Failed to delete domain"
//	@Failure		404			{object}	string	"Domain not found"
//	@Router			/api/v1/domain/{id} [delete]
func (h *Handler) DeleteDomain(c *fiber.Ctx) error {
	var req pb.DeleteDomainRequest

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.DomainId = int64(domainID)

	_, err = h.domainService.DeleteDomain(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrDomainNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrDeleteDomain, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
