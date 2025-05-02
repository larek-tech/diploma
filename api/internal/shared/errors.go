package shared

import (
	"errors"
)

// 401
var (
	// ErrUnauthorized is an error when user failed authorization check.
	ErrUnauthorized = errors.New("unauthorized")
)

// 422
var (
	// ErrInvalidBody is an error when provided an invalid request body that can't be parsed.
	ErrInvalidBody = errors.New("can't parse invalid request body")
)
