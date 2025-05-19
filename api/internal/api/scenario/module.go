package scenario

import (
	"github.com/gofiber/fiber/v2"
)

type scenarioHandler interface {
	CreateScenario(c *fiber.Ctx) error
	GetScenario(c *fiber.Ctx) error
	UpdateScenario(c *fiber.Ctx) error
	DeleteScenario(c *fiber.Ctx) error
	ListScenarios(c *fiber.Ctx) error
	ListScenariosByDomain(c *fiber.Ctx) error
}

// SetupRoutes maps scenario routes.
func SetupRoutes(api fiber.Router, h scenarioHandler) {
	api.Post("/", h.CreateScenario)
	api.Get("/list", h.ListScenarios)
	api.Get("/list_by_domain/:id", h.ListScenariosByDomain)
	api.Get("/:id", h.GetScenario)
	api.Put("/:id", h.UpdateScenario)
	api.Delete("/:id", h.DeleteScenario)
}
