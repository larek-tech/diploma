package role

import (
	"github.com/gofiber/fiber/v2"
)

type roleHandler interface {
	CreateRole(c *fiber.Ctx) error
	GetRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error
	ListRoles(c *fiber.Ctx) error
	SetRole(c *fiber.Ctx) error
	RemoveRole(c *fiber.Ctx) error
}

// SetupRoutes map role routes.
func SetupRoutes(api fiber.Router, h roleHandler) {
	api.Post("/", h.CreateRole)
	api.Get("/list", h.ListRoles)
	api.Get("/:id", h.GetRole)
	api.Put("/set", h.SetRole)
	api.Put("/remove", h.RemoveRole)
	api.Delete("/:id", h.DeleteRole)
}
