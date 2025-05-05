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

// UpdatePermittedUsers godoc
//
//	@Summary		Update permitted users.
//	@Description	Updates list of users permitted to domain.
//	@Tags			domain
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			domainID	path		int					true	"Requested domain ID"
//	@Param			req			body		pb.PermittedUsers	true	"New list of permitted users"
//	@Success		200			{object}	pb.PermittedUsers	"Updated users permissions"
//	@Failure		404			{object}	string				"Domain not found"
//	@Router			/api/v1/domain/permissions/users/{id} [put]
func (h *Handler) UpdatePermittedUsers(c *fiber.Ctx) error {
	var req pb.PermittedUsers
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	domainID, err := c.ParamsInt(domainIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ResourceId = int64(domainID)

	userID, ok := c.Locals(shared.UserIDKey).(int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user ID in context")
	}
	userRoleIDs, ok := c.Locals(shared.UserRolesKey).([]int64)
	if !ok {
		return errs.WrapErr(shared.ErrUnauthorized, "no user users in context")
	}
	ctx := auth.PushUserMeta(c.UserContext(), &authpb.UserAuthMetadata{
		UserId: userID,
		Roles:  userRoleIDs,
	})

	resp, err := h.domainService.UpdatePermittedUsers(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
