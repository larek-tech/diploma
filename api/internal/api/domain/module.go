package domain

import (
	"github.com/gofiber/fiber/v2"
)

type domainHandler interface {
	CreateDomain(c *fiber.Ctx) error
	GetDomain(c *fiber.Ctx) error
	UpdateDomain(c *fiber.Ctx) error
	DeleteDomain(c *fiber.Ctx) error
	ListDomains(c *fiber.Ctx) error
	GetPermittedUsers(c *fiber.Ctx) error
	GetPermittedRoles(c *fiber.Ctx) error
	UpdatePermittedUsers(c *fiber.Ctx) error
	UpdatePermittedRoles(c *fiber.Ctx) error
}

// SetupRoutes maps domain routes.
func SetupRoutes(api fiber.Router, h domainHandler) {
	api.Post("/", h.CreateDomain)
	api.Get("/list", h.ListDomains)
	api.Get("/permissions/users/:id", h.GetPermittedUsers)
	api.Get("/permissions/roles/:id", h.GetPermittedRoles)
	api.Put("/permissions/users/:id", h.UpdatePermittedUsers)
	api.Put("/permissions/roles/:id", h.UpdatePermittedRoles)
	api.Get("/:id", h.GetDomain)
	api.Put("/:id", h.UpdateDomain)
	api.Delete("/:id", h.DeleteDomain)
}
