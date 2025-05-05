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

// UpdatePermittedRoles godoc
//
//	@Summary		Update permitted roles.
//	@Description	Updates list of roles permitted to source.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			sourceID	path		int					true	"Requested source ID"
//	@Param			req			body		pb.PermittedRoles	true	"New list of permitted roles"
//	@Success		200			{object}	pb.PermittedRoles	"Updated roles permissions"
//	@Failure		404			{object}	string				"Source not found"
//	@Router			/api/v1/source/permissions/roles/{id} [put]
func (h *Handler) UpdatePermittedRoles(c *fiber.Ctx) error {
	var req pb.PermittedRoles
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	sourceID, err := c.ParamsInt(sourceIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.ResourceId = int64(sourceID)

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

	resp, err := h.sourceService.UpdatePermittedRoles(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
