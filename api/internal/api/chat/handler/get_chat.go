package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetChat godoc
//
//	@Summary		Get chat.
//	@Description	Returns chat with messages within it.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string	yes	"Chat ID"
//	@Success		200	{object}	pb.Chat	"Returned chat"
//	@Failure		400	{object}	string	"Failed to get chat"
//	@Failure		404	{object}	string	"Chat not found"
//	@Router			/api/v1/chat/history/{id} [get]
func (h *Handler) GetChat(c *fiber.Ctx) error {
	chatID := c.Params(chatIDParam)

	req := &pb.GetChatRequest{ChatId: chatID}
	chat, err := h.chatService.GetChat(c.UserContext(), req)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrChatNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrGetChat, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(chat)
}
