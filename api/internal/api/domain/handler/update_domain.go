package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/auth"
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
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
//	@Param			domainID	path		int						true	"Domain ID"
//	@Param			req			body		pb.UpdateDomainRequest	true	"Update params"
//	@Success		200			{object}	pb.Domain				"Domain updated"
//	@Failure		400			{object}	string					"Failed to update domain"
//	@Failure		404			{object}	string					"Domain not found"
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

	resp, err := h.domainService.UpdateDomain(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrDomainNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrUpdateDomain, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
