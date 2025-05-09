package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateChat godoc
//
//	@Summary		Create new chat.
//	@Description	Creates new chat for RAG system.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		201	{object}	pb.Chat	"Chat successfully created"
//	@Failure		400	{object}	string	"Failed to create chat"
//	@Router			/api/v1/chat/ [post]
func (h *Handler) CreateChat(c *fiber.Ctx) error {
	chat, err := h.chatService.CreateChat(c.UserContext(), &emptypb.Empty{})
	if err != nil {
		return errs.WrapErr(shared.ErrCreateChat, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(chat)
}
