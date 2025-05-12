package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CancelQuery godoc
//
//	@Summary		Cancel query.
//	@Description	Cancel processing of query (all dependant jobs) by id.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int		yes	"Query ID"
//	@Success		204	{object}	pb.Chat	"Processing query successfully canceled"
//	@Failure		400	{object}	string	"Failed to cancel query"
//	@Failure		403	{object}	string	"No access to cancel query"
//	@Router			/api/v1/chat/cancel/{id} [post]
func (h *Handler) CancelQuery(c *fiber.Ctx) error {
	queryID, err := c.ParamsInt(queryIDParam)
	if err != nil {
		return errs.WrapErr(shared.ErrInvalidParams)
	}

	req := &pb.CancelProcessingRequest{
		QueryId: int64(queryID),
	}
	_, err = h.chatService.CancelProcessing(c.UserContext(), req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrChatNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrCancelQuery, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
