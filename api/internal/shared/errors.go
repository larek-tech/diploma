package shared

import (
	"errors"
)

// 400
var (
	// ErrCreateSource is an error when failed to create domain source.
	ErrCreateSource = errors.New("failed to create source")
	// ErrGetSource is an error when failed to get domain source.
	ErrGetSource = errors.New("failed to get source")
	// ErrUpdateSource is an error when failed to update domain source.
	ErrUpdateSource = errors.New("failed to update source")
	// ErrDeleteSource is an error when failed to delete domain source.
	ErrDeleteSource = errors.New("failed to delete source")
	// ErrListSources is an error when failed to list domain sources.
	ErrListSources = errors.New("failed to list sources")
)

// 401
var (
	// ErrUnauthorized is an error when user failed authorization check.
	ErrUnauthorized = errors.New("unauthorized")
)

// 404
var (
	// ErrSourceNotFound is an error when no source was found with requested ID.
	ErrSourceNotFound = errors.New("no source with such ID")
)

// 422
var (
	// ErrInvalidBody is an error when provided an invalid request body that can't be parsed.
	ErrInvalidBody = errors.New("can't parse invalid request body")
	// ErrInvalidParams is an error when provided invalid path or query param that can't be parsed.
	ErrInvalidParams = errors.New("can't parse invalid path or query params")
)
