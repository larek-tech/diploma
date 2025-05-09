package handler

import (
	"github.com/larek-tech/diploma/api/internal/chat/pb"
)

// Handler implements chat methods on transport level.
type Handler struct {
	chatService pb.ChatServiceClient
}

// New creates new Handler.
func New(chatService pb.ChatServiceClient) *Handler {
	return &Handler{chatService: chatService}
}
