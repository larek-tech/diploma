package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteChat godoc
//
//	@Summary		Delete chat.
//	@Description	Soft delete chat by ID.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			chatID	path		string	yes	"Chat ID"
//	@Success		204		{object}	pb.Chat	"Chat successfully deleted"
//	@Failure		400		{object}	string	"Failed to delete chat"
//	@Failure		404		{object}	string	"Chat not found"
//	@Router			/api/v1/chat/{id} [delete]
func (h *Handler) DeleteChat(c *fiber.Ctx) error {
	chatID := c.Params(chatIDParam)

	req := &pb.DeleteChatRequest{ChatId: chatID}
	chat, err := h.chatService.DeleteChat(c.UserContext(), req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrChatNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrDeleteChat, err.Error())
	}

	return c.Status(fiber.StatusNoContent).JSON(chat)
}
