// Package response provides a way to handle errors and return them as JSON responses via fiber framework.
package response

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/yogenyslav/pkg/errs"
)

// ErrorResponse is a struct that holds the error message and status code.
type ErrorResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"-"`
}

// ErrorHandler is a struct that holds the error status map.
type ErrorHandler struct {
	status map[error]ErrorResponse
}

// NewErrorHandler creates a new ErrorHandler instance with the given error status map.
func NewErrorHandler(errStatus map[error]ErrorResponse) ErrorHandler {
	status := map[error]ErrorResponse{
		pgx.ErrNoRows: {
			Msg:    "no rows found",
			Status: http.StatusNotFound,
		},
		fiber.ErrUnprocessableEntity: {
			Msg:    "validation error",
			Status: http.StatusUnprocessableEntity,
		},
	}

	for k, v := range errStatus {
		status[k] = v
	}

	return ErrorHandler{
		status: status,
	}
}

// Handler is a method that handles the error and returns a JSON response.
// Should be used as a fiber.Config.ErrorHandler.
func (h ErrorHandler) Handler(logger *zerolog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		e := h.getErrorResponse(err)
		logger.Err(err).Msg(e.Msg)
		return c.Status(e.Status).JSON(e) //nolint:wrapcheck // no need to wrap
	}
}

func (h ErrorHandler) getErrorResponse(err error) ErrorResponse {
	var (
		ok bool
		e  ErrorResponse
	)

	if errs.CheckPageNotFound(err) {
		return ErrorResponse{
			Msg:    "page not found",
			Status: http.StatusNotFound,
		}
	}

	if errs.CheckDuplicateKey(err) {
		return ErrorResponse{
			Msg:    "duplicate key",
			Status: http.StatusBadRequest,
		}
	}

	if errs.CheckValidationError(err) {
		return ErrorResponse{
			Msg:    err.Error(),
			Status: http.StatusUnprocessableEntity,
		}
	}

	for k, v := range h.status {
		if errors.Is(err, k) {
			ok = true
			e = v
			break
		}
	}

	if !ok {
		e = ErrorResponse{
			Msg:    "unknown error",
			Status: http.StatusInternalServerError,
		}
	}
	if e.Msg == "" {
		e.Msg = err.Error()
	}

	return e
}
