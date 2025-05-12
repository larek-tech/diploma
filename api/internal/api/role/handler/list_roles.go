package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// ListRoles godoc
//
//	@Summary		List roles.
//	@Description	List roles.
//	@Tags			role
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			offset	query		uint					true	"Pagination offset"
//	@Param			limit	query		uint					true	"Pagination limit"
//	@Success		200		{object}	pb.ListRolesResponse	"List of roles"
//	@Failure		400		{object}	string					"Failed to list roles"
//	@Router			/api/v1/role/list [get]
func (h *Handler) ListRoles(c *fiber.Ctx) error {
	offset := c.QueryInt(offsetParam, 0)
	limit := c.QueryInt(limitParam, 10)
	if offset < 0 || limit < 0 {
		return errs.WrapErr(shared.ErrInvalidParams, fmt.Sprintf("offset=%d, limit=%d", offset, limit))
	}

	req := pb.ListRolesRequest{
		Offset: uint64(offset),
		Limit:  uint64(limit),
	}
	resp, err := h.roleService.ListRoles(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrListRoles, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
