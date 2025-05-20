package result_message

import (
	"context"
	"log/slog"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h Handler) Handle(ctx context.Context, job any) error {
	slog.Debug("handled result message", "job", job)
	return nil
}
