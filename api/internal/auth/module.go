package auth

import (
	"github.com/gofiber/fiber/v2"
)

type authHandler interface {
	Login(c *fiber.Ctx) error
}

// SetupRoutes maps auth routes.
func SetupRoutes(auth fiber.Router, h authHandler) {
	auth.Post("/login", h.Login)
}
