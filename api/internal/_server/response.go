package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/response"
)

var (
	errMap = map[error]response.ErrorResponse{
		// 400
		shared.ErrCreateSource: {
			Msg:    "failed creating source",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrGetSource: {
			Msg:    "failed getting source",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrUpdateSource: {
			Msg:    "failed updating source",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrDeleteSource: {
			Msg:    "failed deleting source",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrListSources: {
			Msg:    "failed listing available sources",
			Status: fiber.StatusBadRequest,
		},
		// 401
		shared.ErrUnauthorized: {
			Msg:    "unauthorized",
			Status: fiber.StatusUnauthorized,
		},
		// 404
		shared.ErrSourceNotFound: {
			Msg:    "no source with such id",
			Status: fiber.StatusNotFound,
		},
		// 422
		shared.ErrInvalidBody: {
			Msg:    "can't parse request body",
			Status: fiber.StatusUnprocessableEntity,
		},
		shared.ErrInvalidParams: {
			Msg:    "can't parse path or query params",
			Status: fiber.StatusUnprocessableEntity,
		},
	}
)
