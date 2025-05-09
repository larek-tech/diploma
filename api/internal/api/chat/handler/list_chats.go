package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// ListChats godoc
//
//	@Summary		List chats.
//	@Description	List chats created by user.
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			offset	query		uint					true	"Pagination offset"
//	@Param			limit	query		uint					true	"Pagination limit"
//	@Success		200		{object}	pb.ListChatsResponse	"List of chats"
//	@Failure		400		{object}	string					"Failed to list chats"
//	@Router			/api/v1/chat/list [get]
func (h *Handler) ListChats(c *fiber.Ctx) error {
	offset := c.QueryInt(offsetParam, 0)
	limit := c.QueryInt(limitParam, 10)
	if offset < 0 || limit < 0 {
		return errs.WrapErr(shared.ErrInvalidParams, fmt.Sprintf("offset=%d, limit=%d", offset, limit))
	}

	req := pb.ListChatsRequest{
		Offset: uint64(offset),
		Limit:  uint64(limit),
	}
	resp, err := h.chatService.ListChats(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrListChats, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
