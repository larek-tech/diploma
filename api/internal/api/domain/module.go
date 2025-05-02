package domain

import (
	"github.com/gofiber/fiber/v2"
)

type domainHandler interface {
	CreateSource(c *fiber.Ctx) error
	GetSource(c *fiber.Ctx) error
	UpdateSource(c *fiber.Ctx) error
	DeleteSource(c *fiber.Ctx) error
	ListSources(c *fiber.Ctx) error
}

func Setup(api fiber.Router, h domainHandler) {

	api.Post("/", h.CreateSource)
	api.Get("/list", h.ListSources)
	api.Put("/", h.CreateSource)
	api.Get("/:id", h.GetSource)
	api.Delete("/:id", h.DeleteSource)
}
