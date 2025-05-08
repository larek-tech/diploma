package controller

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/chat/internal/chat/model"
	"github.com/larek-tech/diploma/chat/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrNoAccessToChat is an error when user can't edit chat.
	ErrNoAccessToChat = errors.New("user has no access to edit chat")
)

type chatRepo interface {
	InsertChat(ctx context.Context, chat model.ChatDao) (uuid.UUID, error)
	InsertQuery(ctx context.Context, q model.QueryDao) (int64, error)
	InsertResponse(ctx context.Context, resp model.ResponseDao) (int64, error)
	GetChat(ctx context.Context, chatID uuid.UUID) (model.ChatDao, error)
	GetChatUserID(ctx context.Context, chatID uuid.UUID) (int64, error)
	GetResponseByID(ctx context.Context, respID int64) (model.ResponseDao, error)
	UpdateChatTitle(ctx context.Context, title string, chatID uuid.UUID) error
	UpdateResponse(ctx context.Context, resp model.ResponseDao) error
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	SoftDeleteChat(ctx context.Context, chatID uuid.UUID) error
	ListChats(ctx context.Context, offset, limit uint64, userID int64) ([]model.ChatDao, error)
}

// Controller implements chat methods on logic layer.
type Controller struct {
	cr        chatRepo
	tracer    trace.Tracer
	mlService pb.MLServiceClient
}

// New creates new Controller.
func New(cr chatRepo, tracer trace.Tracer, mlService pb.MLServiceClient) *Controller {
	return &Controller{
		cr:        cr,
		tracer:    tracer,
		mlService: mlService,
	}
}
