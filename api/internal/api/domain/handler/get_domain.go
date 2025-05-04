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

// GetDomain godoc
//
//	@Summary		Get domain.
//	@Description	Returns information about domain.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			domainID	path		int						true	"Requested domain ID"
//	@Success		200			{object}	pb.GetDomainResponse	"Domain"
//	@Failure		400			{object}	string					"Failed to get domain"
//	@Failure		404			{object}	string					"Domain not found"
//	@Router			/api/v1/domain/{id} [get]
func (h *Handler) GetDomain(c *fiber.Ctx) error {
	var req pb.GetDomainRequest

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

	resp, err := h.domainService.GetDomain(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrDomainNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrGetDomain, err.Error())
	}

	if resp.GetDomain() == nil {
		return errs.WrapErr(shared.ErrDomainNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
