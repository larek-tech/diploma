package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/larek-tech/diploma/api/internal/api/chat/model"
	"github.com/larek-tech/diploma/api/internal/auth"
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/chat/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
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
	log.Err(errs.WrapErr(err)).Msg("chat error")
	sendMsg(c, &model.SocketMessage{
		Type:   model.TypeError,
		IsLast: true,
		Err:    errors.New(msg),
	})
}

func sendMsg(c *websocket.Conn, msg *model.SocketMessage) {
	if err := c.WriteJSON(msg); err != nil {
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

	return userMeta.GetMeta(), nil
}

// Chat handles websocket connection for sending messages.
func (h *Handler) Chat(c *websocket.Conn) {
	ctx, span := h.tracer.Start(context.Background(), "Handler.Chat")
	defer span.End()

	c.SetCloseHandler(closeHandler)
	defer func() {
		if e := c.Close(); e != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("close websocket conn")
		}
	}()

	chatID := c.Params(chatIDParam)
	addr := c.LocalAddr().String()
	span.SetAttributes(
		attribute.String("chatID", chatID),
		attribute.String("addr", addr),
	)
	log.Info().Str("addr", addr).Msg("new conn")

	userMeta, err := h.authorize(c, ctx)
	if err != nil {
		sendErr(c, errs.WrapErr(err), "unauthorized")
		return
	}

	span.SetAttributes(attribute.Int64("userID", userMeta.GetUserId()))

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
		chunk      *pb.ChunkedResponse
	)

	for {
		msg, err = getMsg(c)
		if err != nil {
			sendErr(c, errs.WrapErr(err), "read next query")
			return
		}

		if msg.Type != model.TypeQuery {
			sendErr(
				c,
				errs.WrapErr(shared.ErrInvalidBody),
				fmt.Sprintf("unexpected message type: got %s, want %s", msg.Type, model.TypeQuery),
			)
			return
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
				return
			}
		}

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
			return
		}

		chunk, err = stream.Recv()
		if err != nil {
			sendErr(c, errs.WrapErr(err), "read next chunk")
			return
		}

		// not last chunk
		if sourceIDs := chunk.GetSourceIds(); sourceIDs == nil {
			msg = model.SocketMessage{
				Type:      model.TypeChunk,
				Content:   chunk.GetContent(),
				IsChunked: true,
			}
		} else {
			msg = model.SocketMessage{
				Type:      model.TypeChunk,
				IsChunked: true,
				IsLast:    true,
				SourceIDs: sourceIDs,
			}
		}
		sendMsg(c, &msg)
	}
}
