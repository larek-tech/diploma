package user

import (
	"github.com/gofiber/fiber/v2"
)

type userHandler interface {
	CreateUser(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	ListUsers(c *fiber.Ctx) error
}

// SetupRoutes map user routes.
func SetupRoutes(api fiber.Router, h userHandler) {
	api.Post("/", h.CreateUser)
	api.Get("/list", h.ListUsers)
	api.Get("/:id", h.GetUser)
	api.Delete("/:id", h.DeleteUser)
}
