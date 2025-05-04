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

// GetSource godoc
//
//	@Summary		Get source.
//	@Description	Returns information about source and its update parameters.
//	@Tags			source
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			sourceID	path		int						true	"Requested source ID"
//	@Success		200			{object}	pb.GetSourceResponse	"Source"
//	@Failure		400			{object}	string					"Failed to get source"
//	@Failure		404			{object}	string					"Source not found"
//	@Router			/api/v1/source/{id} [get]
func (h *Handler) GetSource(c *fiber.Ctx) error {
	var req pb.GetSourceRequest

	sourceID, err := c.ParamsInt(sourceIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams, err.Error())
	}
	req.SourceId = int64(sourceID)

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

	resp, err := h.sourceService.GetSource(ctx, &req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrSourceNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrGetSource, err.Error())
	}

	if resp.GetSource() == nil {
		return errs.WrapErr(shared.ErrSourceNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
