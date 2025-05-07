package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// ListDomains godoc
//
//	@Summary		List domains.
//	@Description	List domains to which user has access.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			offset	query		uint					true	"Pagination offset"
//	@Param			limit	query		uint					true	"Pagination limit"
//	@Success		200		{object}	pb.ListDomainsResponse	"List of domains"
//	@Failure		400		{object}	string					"Failed to list domains"
//	@Router			/api/v1/domain/list [get]
func (h *Handler) ListDomains(c *fiber.Ctx) error {
	offset := c.QueryInt(offsetParam, 0)
	limit := c.QueryInt(limitParam, 10)
	if offset < 0 || limit < 0 {
		return errs.WrapErr(shared.ErrInvalidParams, fmt.Sprintf("offset=%d, limit=%d", offset, limit))
	}

	req := pb.ListDomainsRequest{
		Offset: uint64(offset),
		Limit:  uint64(limit),
	}
	resp, err := h.domainService.ListDomains(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrListDomains, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
