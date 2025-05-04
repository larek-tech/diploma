package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/auth"
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// ListSources godoc
//
//	@Summary		List sources.
//	@Description	List sources to which user has access.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			offset	query		uint					true	"Pagination offset"
//	@Param			limit	query		uint					true	"Pagination limit"
//	@Success		200		{object}	pb.ListSourcesResponse	"List of sources"
//	@Failure		400		{object}	string					"Failed to list sources"
//	@Router			/api/v1/source/list [get]
func (h *Handler) ListSources(c *fiber.Ctx) error {
	offset := c.QueryInt(offsetParam, 0)
	limit := c.QueryInt(limitParam, 10)
	if offset < 0 || limit < 0 {
		return errs.WrapErr(shared.ErrInvalidParams, fmt.Sprintf("offset=%d, limit=%d", offset, limit))
	}

	userID, ok := c.Locals(shared.UserIDKey).(int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user ID in context")
	}
	userRoleIDs, ok := c.Locals(shared.UserRolesKey).([]int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user roles in context")
	}
	ctx := auth.PushUserMeta(c.UserContext(), &authpb.UserAuthMetadata{
		UserId: userID,
		Roles:  userRoleIDs,
	})

	req := pb.ListSourcesRequest{
		Offset: uint64(offset),
		Limit:  uint64(limit),
	}
	resp, err := h.sourceService.ListSources(ctx, &req)
	if err != nil {
		return errs.WrapErr(shared.ErrListSources, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
