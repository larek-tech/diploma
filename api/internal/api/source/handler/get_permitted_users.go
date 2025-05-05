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

// GetPermittedUsers godoc
//
//	@Summary		Get permitted users.
//	@Description	Returns list of users permitted to source.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			sourceID	path		int						true	"Requested source ID"
//	@Success		200			{object}	pb.GetSourceResponse	"Permitted users"
//	@Failure		404			{object}	string					"Source not found"
//	@Router			/api/v1/source/permissions/users/{id} [get]
func (h *Handler) GetPermittedUsers(c *fiber.Ctx) error {
	var req pb.GetResourcePermissionsRequest

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

	resp, err := h.sourceService.GetPermittedUsers(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		return errs.WrapErr(errs.WrapErr(err))
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
