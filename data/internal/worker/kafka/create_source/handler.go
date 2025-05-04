package create_source

import (
	"context"
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h Handler) Handle(ctx context.Context, msg *kafka.Message) error {
	slog.Info("received new msg", "msg", msg)
	return nil
}
