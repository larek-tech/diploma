package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// ListSourcesByDomain godoc
//
//	@Summary		List sources by domain.
//	@Description	List sources by domain to which user has access.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			offset	query		uint					true	"Pagination offset"
//	@Param			limit	query		uint					true	"Pagination limit"
//	@Param			id		path		int64					true	"Domain ID"
//	@Success		200		{object}	pb.ListSourcesResponse	"List of sources by domain"
//	@Failure		400		{object}	string					"Failed to list sources by domain"
//	@Router			/api/v1/source/list_by_domain/{id} [get]
func (h *Handler) ListSourcesByDomain(c *fiber.Ctx) error {
	offset := c.QueryInt(offsetParam, 0)
	limit := c.QueryInt(limitParam, 10)
	if offset < 0 || limit < 0 {
		return errs.WrapErr(shared.ErrInvalidParams, fmt.Sprintf("offset=%d, limit=%d", offset, limit))
	}

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, fmt.Sprintf("invalid domain ID: %s", err.Error()))
	}

	req := pb.ListSourcesByDomainRequest{
		DomainId: int64(domainID),
		Offset:   uint64(offset),
		Limit:    uint64(limit),
	}
	resp, err := h.sourceService.ListSourcesByDomain(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrListSources, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
