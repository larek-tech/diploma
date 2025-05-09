package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RenameChat godoc
//
//	@Summary		Update chat.
//	@Description	Update chat information.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			chatID	path		int						true	"Chat ID"
//	@Param			req		body		pb.RenameChatRequest	true	"Update params"
//	@Success		200		{object}	pb.Chat					"Chat updated"
//	@Failure		400		{object}	string					"Failed to update chat"
//	@Failure		404		{object}	string					"Chat not found"
//	@Router			/api/v1/chat/{id} [put]
func (h *Handler) RenameChat(c *fiber.Ctx) error {
	var req pb.RenameChatRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	chatID := c.Params(chatIDParam)
	req.ChatId = chatID

	resp, err := h.chatService.RenameChat(c.UserContext(), &req)
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return errs.WrapErr(shared.ErrForbidden, err.Error())
		}
		if status.Code(err) == codes.NotFound {
			return errs.WrapErr(shared.ErrChatNotFound, err.Error())
		}
		return errs.WrapErr(shared.ErrUpdateChat, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
