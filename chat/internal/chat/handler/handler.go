package handler

import (
	"context"
	"errors"

	authpb "github.com/larek-tech/diploma/chat/internal/auth/pb"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"go.opentelemetry.io/otel/trace"
)

var (
	// ErrProcessQueryTimeout is an error when processing query breaks timeout.
	ErrProcessQueryTimeout = errors.New("process query timeout")
	// ErrEmptySourceIDs is an error when list of source ids is empty.
	ErrEmptySourceIDs = errors.New("empty source ids list")
)

type chatController interface {
	CreateChat(ctx context.Context, userID int64) (*pb.Chat, error)
	GetChat(ctx context.Context, chatID string) (*pb.Chat, error)
	RenameChat(ctx context.Context, req *pb.RenameChatRequest, meta *authpb.UserAuthMetadata) (*pb.Chat, error)
	DeleteChat(ctx context.Context, chatID string, meta *authpb.UserAuthMetadata) error
	CleanupChat(ctx context.Context, chatID string) error
	ListChats(ctx context.Context, req *pb.ListChatsRequest, meta *authpb.UserAuthMetadata) (*pb.ListChatsResponse, error)
	ProcessQuery(ctx context.Context, req *pb.ProcessQueryRequest, out chan *pb.ChunkedResponse, errCh chan error)
	CancelProcessing(ctx context.Context, req *pb.CancelProcessingRequest, meta *authpb.UserAuthMetadata) error
}

// Handler implements chat methods on transport level.
type Handler struct {
	pb.UnimplementedChatServiceServer
	cc     chatController
	tracer trace.Tracer
}

// New creates new Handler.
func New(cc chatController, tracer trace.Tracer) *Handler {
	return &Handler{
		cc:     cc,
		tracer: tracer,
	}
}
