package result_message

import (
	"context"
	"log/slog"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas/messages"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h Handler) Handle(ctx context.Context, result messages.ResultMessage) error {
	slog.Info("handled result message", "result", result)
	return nil
}
