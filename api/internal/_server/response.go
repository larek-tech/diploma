package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/response"
)

var (
	errMap = map[error]response.ErrorResponse{
		shared.ErrUnauthorized: {
			Status: fiber.StatusUnauthorized,
		},
	}
)
