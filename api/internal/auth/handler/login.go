package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// Login authorizes user with credentials.
func (h *Handler) Login(c *fiber.Ctx) error {
	var req pb.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errs.WrapErr(shared.ErrInvalidBody, err.Error())
	}

	resp, err := h.authService.Login(c.UserContext(), &req)
	if err != nil {
		return errs.WrapErr(shared.ErrUnauthorized, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
