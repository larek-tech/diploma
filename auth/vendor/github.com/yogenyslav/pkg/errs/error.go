// Package errs provides common errors and method to handle them.
package errs

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	// ErrDuplicateKey is an error for postgres unique key violation.
	ErrDuplicateKey = errors.New("duplicate key")
	// ErrPageNotFound for page not found handlers.
	ErrPageNotFound = errors.New("page not found")
	// ErrValidation for handling validation errors.
	ErrValidation = errors.New("validation error")
)

// CheckDuplicateKey checks if the error is a postgres duplicate key violation.
func CheckDuplicateKey(err error) bool {
	var pgError *pgconn.PgError
	return errors.As(err, &pgError) && pgError.Code == "23505"
}

// CheckPageNotFound checks if the error is a fiber page not found error.
func CheckPageNotFound(err error) bool {
	var fiberError *fiber.Error
	return errors.As(err, &fiberError) && fiberError.Code == http.StatusNotFound
}

// CheckValidationError checks if the error is a validation error.
func CheckValidationError(err error) bool {
	var validationErrors validator.ValidationErrors
	return errors.As(err, &validationErrors)
}

// WrapErr fits the error in a chain, reports source file and provides optional description.
func WrapErr(e error, desc ...string) error {
	if e == nil {
		return nil
	}
	var d string
	if len(desc) > 0 {
		d = desc[0]
	}
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("undefined call %s -> %w", d, e)
	}
	return fmt.Errorf("%s:%d %s -> %w", file, line, d, e)
}
