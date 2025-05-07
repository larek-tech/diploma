package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// CreateDomain godoc
//
//	@Summary		Create new domain.
//	@Description	Creates new domain (group of sources used for RAG vector search).
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			req	body		pb.CreateDomainRequest	true	"Input data for creating domain"
//	@Success		201	{object}	pb.Domain				"Domain successfully created"
//	@Failure		400	{object}	string					"Failed to create domain"
//	@Router			/api/v1/domain/ [post]
func (h *Handler) CreateDomain(c *fiber.Ctx) error {
	var req pb.CreateDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.domainService.CreateDomain(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrCreateDomain, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}
