package chat

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

type chatHandler interface {
	CreateChat(c *fiber.Ctx) error
	GetChat(c *fiber.Ctx) error
	RenameChat(c *fiber.Ctx) error
	DeleteChat(c *fiber.Ctx) error
	ListChats(c *fiber.Ctx) error
	CancelQuery(c *fiber.Ctx) error
	Chat(c *websocket.Conn)
}

// SetupRoutes map chat routes.
func SetupRoutes(api fiber.Router, h chatHandler, wsConfig websocket.Config) {
	api.Post("/", h.CreateChat)
	api.Get("/list", h.ListChats)
	api.Get("/history/:id", h.GetChat)
	api.Put("/:id", h.RenameChat)
	api.Delete("/:id", h.DeleteChat)

	api.Route("/:id", func(ws fiber.Router) {
		ws.Use("/:id", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}
			return errs.WrapErr(shared.ErrWsProtocolRequired)
		})
		ws.Get("/:id", websocket.New(h.Chat, wsConfig))
	}, "ws")
}
