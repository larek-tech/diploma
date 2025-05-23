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
	domainpb "github.com/larek-tech/diploma/api/internal/domain/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	history, err := h.chatService.GetChat(ctx, &pb.GetChatRequest{ChatId: chatID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			sendErr(c, errs.WrapErr(err), "chat not found")
			return
		}
		sendErr(c, errs.WrapErr(err), "get chat error")
		return
	}

	var (
		msg              model.SocketMessage
		titleResp        *domainpb.ProcessFirstQueryResponse
		processReq       *pb.ProcessQueryRequest
		chunk            *model.SocketMessage
		scenario         *domainpb.Scenario
		domain           *domainpb.Domain
		sourceIDs        *domainpb.GetSourceIDsResponse
		scenarioMetadata []byte
		errStream        error
		firstMessage     bool
	)

	if len(history.GetContent()) == 0 {
		firstMessage = true
	}

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

		ctx = context.WithValue(ctx, "chat-id", chatID)

		scenario, err = h.scenarioService.GetScenario(ctx, &domainpb.GetScenarioRequest{ScenarioId: msg.ScenarioID})
		if err != nil {
			sendErr(c, errs.WrapErr(err), "get scenario by id")
			continue
		}

		scenarioMetadata, err = json.Marshal(scenario)
		if err != nil {
			sendErr(c, errs.WrapErr(err), "process scenario")
			continue
		}

		domain, err = h.domainService.GetDomain(ctx, &domainpb.GetDomainRequest{DomainId: msg.DomainID})
		if err != nil {
			sendErr(c, errs.WrapErr(err), "get domain by id")
			continue
		}

		sourceIDs, err = h.sourceService.GetSourceIDs(ctx, &domainpb.GetSourceIDsRequest{SourceIds: domain.GetSourceIds()})
		if err != nil {
			sendErr(c, errs.WrapErr(err), "get source ids")
			continue
		}

		processReq = &pb.ProcessQueryRequest{
			UserId:    userMeta.GetUserId(),
			ChatId:    chatID,
			Content:   msg.Content,
			DomainId:  msg.DomainID,
			Scenario:  scenarioMetadata,
			SourceIds: sourceIDs.GetSourceIds(),
		}
		stream, e := h.chatService.ProcessQuery(ctx, processReq)
		if e != nil {
			sendErr(c, errs.WrapErr(e), "start processing query")
			continue
		}

		if firstMessage {
			titleResp, err = h.mlService.ProcessFirstQuery(ctx, &domainpb.ProcessFirstQueryRequest{
				Query: msg.Content,
			})
			if err != nil {
				sendErr(c, errs.WrapErr(err), "summarize first message")
			} else {
				_, err = h.chatService.RenameChat(ctx, &pb.RenameChatRequest{
					ChatId: chatID,
					Title:  titleResp.GetQuery(),
				})
				if err != nil {
					sendErr(c, errs.WrapErr(err), "update chat title")
				}
			}

			firstMessage = false
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
		}
		sendMsg(c, chunk)
	}
}
