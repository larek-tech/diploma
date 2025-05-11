package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gofiber/contrib/websocket"
	"github.com/larek-tech/diploma/api/internal/api/chat/model"
	"github.com/larek-tech/diploma/api/internal/auth"
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
)

func closeHandler(code int, text string) error {
	log.Info().Int("code", code).Str("text", text).Msg("close websocket conn")
	return nil
}

func getMsg(c *websocket.Conn) (model.SocketMessage, error) {
	var msg model.SocketMessage

	if err := c.ReadJSON(&msg); err != nil {
		return msg, errs.WrapErr(err, "read message")
	}

	return msg, nil
}

func sendErr(c *websocket.Conn, err error, msg string) {
	if websocket.IsCloseError(err, websocket.CloseGoingAway) {
		return
	}

	log.Err(errs.WrapErr(err)).Msg("chat error")
	sendMsg(c, &model.SocketMessage{
		Type:   model.TypeError,
		IsLast: true,
		Err:    msg,
	})
}

func sendMsg(c *websocket.Conn, msg *model.SocketMessage) {
	log.Debug().Any("resp", *msg).Msg("send message")
	if err := c.WriteJSON(*msg); err != nil {
		log.Warn().Err(errs.WrapErr(err)).Msg("send message")
	}
}

func (h *Handler) authorize(c *websocket.Conn, ctx context.Context) (*authpb.UserAuthMetadata, error) {
	ctx, span := h.tracer.Start(ctx, "Handler.authorize")
	defer span.End()

	credentials, err := getMsg(c)
	if err != nil {
		return nil, errs.WrapErr(err, "get auth credentials")
	}

	if credentials.Type != model.TypeAuth {
		return nil, errs.WrapErr(
			shared.ErrUnauthorized,
			fmt.Sprintf("unexpected message type: got %s, want %s", credentials.Type, model.TypeAuth),
		)
	}

	validateReq := &authpb.ValidateRequest{Token: credentials.Content}
	userMeta, err := h.authService.Validate(ctx, validateReq)
	if err != nil {
		return nil, errs.WrapErr(err, "validate token")
	}

	span.SetAttributes(attribute.Int64("userID", userMeta.GetMeta().GetUserId()))

	return userMeta.GetMeta(), nil
}

func (h *Handler) receiveChunk(stream grpc.ServerStreamingClient[pb.ChunkedResponse]) (*model.SocketMessage, error) {
	chunk, err := stream.Recv()
	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, errs.WrapErr(err, "receive next chunk")
	}

	msg := model.SocketMessage{
		Type:      model.TypeChunk,
		Content:   chunk.GetContent(),
		IsChunked: true,
	}
	log.Debug().Str("content", msg.Content).Msg("got chunk")
	return &msg, nil
}

// Chat handles websocket connection for sending messages.
func (h *Handler) Chat(c *websocket.Conn) {
	ctx := context.Background()

	c.SetCloseHandler(closeHandler)
	defer func() {
		if e := c.Close(); e != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("close websocket conn")
		}
	}()

	chatID := c.Params(chatIDParam)
	log.Info().Str("addr", c.LocalAddr().String()).Msg("new conn")

	userMeta, err := h.authorize(c, ctx)
	if err != nil {
		sendErr(c, errs.WrapErr(err), "unauthorized")
		return
	}

	ctx = auth.PushUserMeta(ctx, userMeta)
	defer func() {
		cleanUpReq := &pb.CleanupChatRequest{ChatId: chatID}
		if _, e := h.chatService.CleanupChat(ctx, cleanUpReq); e != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("cleanup chat")
		}
	}()

	var (
		msg        model.SocketMessage
		processReq *pb.ProcessQueryRequest
		chunk      *model.SocketMessage
		errStream  error
	)

	for {
		msg, err = getMsg(c)
		if err != nil {
			sendErr(c, errs.WrapErr(err), "read next query")
			continue
		}

		if msg.Type != model.TypeQuery {
			sendErr(
				c,
				errs.WrapErr(shared.ErrInvalidBody),
				fmt.Sprintf("unexpected message type: got %s, want %s", msg.Type, model.TypeQuery),
			)
			continue
		}

		var (
			scenarioID   *int64 = nil
			scenarioMeta []byte = nil
		)
		if msg.QueryMetadata.Scenario != nil {
			scenarioID = &msg.QueryMetadata.Scenario.Id
			scenarioMeta, err = json.Marshal(msg.QueryMetadata.Scenario)
			if err != nil {
				sendErr(c, errs.WrapErr(shared.ErrInvalidBody), "invalid scenario")
				continue
			}
		}

		ctx = context.WithValue(ctx, "chat-id", chatID)

		processReq = &pb.ProcessQueryRequest{
			UserId:     userMeta.GetUserId(),
			ChatId:     chatID,
			Content:    msg.Content,
			DomainId:   msg.QueryMetadata.DomainID,
			SourceIds:  msg.SourceIDs,
			ScenarioId: scenarioID,
			Metadata:   scenarioMeta,
		}
		stream, e := h.chatService.ProcessQuery(ctx, processReq)
		if e != nil {
			sendErr(c, errs.WrapErr(e), "start processing query")
			continue
		}

		chunk, errStream = h.receiveChunk(stream)
		for chunk != nil && errStream == nil {
			sendMsg(c, chunk)
			chunk, errStream = h.receiveChunk(stream)
		}

		if errStream != nil {
			sendErr(c, errs.WrapErr(err), "receive next chunk")
			continue
		}

		chunk = &model.SocketMessage{
			Type:      model.TypeChunk,
			IsChunked: true,
			IsLast:    true,
			SourceIDs: msg.SourceIDs,
		}
		sendMsg(c, chunk)
	}
}
