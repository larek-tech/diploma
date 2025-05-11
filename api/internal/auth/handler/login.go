package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/errs"
)

// Login godoc
//
//	@Summary		Login user.
//	@Summary		Login user.
//	@Description	Authorizes user with provided credentials.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			req	body		pb.LoginRequest		true	"User credentials"
//	@Success		200	{object}	pb.LoginResponse	"Auth token and metadata"
//	@Failure		401	{object}	string				"Unauthorized"
//	@Router			/auth/v1/login [post]
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
