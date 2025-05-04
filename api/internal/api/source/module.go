package source

import (
	"github.com/gofiber/fiber/v2"
)

type sourceHandler interface {
	CreateSource(c *fiber.Ctx) error
	GetSource(c *fiber.Ctx) error
	UpdateSource(c *fiber.Ctx) error
	DeleteSource(c *fiber.Ctx) error
	ListSources(c *fiber.Ctx) error
}

// SetupRoutes maps source routes.
func SetupRoutes(api fiber.Router, h sourceHandler) {
	api.Post("/", h.CreateSource)
	api.Get("/list", h.ListSources)
	api.Get("/:id", h.GetSource)
	api.Put("/:id", h.UpdateSource)
	api.Delete("/:id", h.DeleteSource)
}
