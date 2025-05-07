package shared

import (
	"errors"
)

// 400
var (
	// ErrCreateSource is an error when failed to create source.
	ErrCreateSource = errors.New("failed to create source")
	// ErrGetSource is an error when failed to get source.
	ErrGetSource = errors.New("failed to get source")
	// ErrUpdateSource is an error when failed to update source.
	ErrUpdateSource = errors.New("failed to update source")
	// ErrDeleteSource is an error when failed to delete source.
	ErrDeleteSource = errors.New("failed to delete source")
	// ErrListSources is an error when failed to list sources.
	ErrListSources = errors.New("failed to list sources")

	// ErrCreateDomain is an error when failed to create domain.
	ErrCreateDomain = errors.New("failed to create domain")
	// ErrGetDomain is an error when failed to get domain.
	ErrGetDomain = errors.New("failed to get domain")
	// ErrUpdateDomain is an error when failed to update domain.
	ErrUpdateDomain = errors.New("failed to update domain")
	// ErrDeleteDomain is an error when failed to delete domain.
	ErrDeleteDomain = errors.New("failed to delete domain")
	// ErrListDomains is an error when failed to list domains.
	ErrListDomains = errors.New("failed to list domains")

	// ErrCreateScenario is an error when failed to create scenario.
	ErrCreateScenario = errors.New("failed to create scenario")
	// ErrGetScenario is an error when failed to get scenario.
	ErrGetScenario = errors.New("failed to get scenario")
	// ErrUpdateScenario is an error when failed to update scenario.
	ErrUpdateScenario = errors.New("failed to update scenario")
	// ErrDeleteScenario is an error when failed to delete scenario.
	ErrDeleteScenario = errors.New("failed to delete scenario")
	// ErrListScenarios is an error when failed to list scenarios.
	ErrListScenarios = errors.New("failed to list scenarios")
)

// 401
var (
	// ErrUnauthorized is an error when user failed authorization check.
	ErrUnauthorized = errors.New("unauthorized")
)

// 403
var (
	// ErrForbidden is an error when user tries to access forbidden resource.
	ErrForbidden = errors.New("forbidden")
)

// 404
var (
	// ErrSourceNotFound is an error when no source was found.
	ErrSourceNotFound = errors.New("source not found")
	// ErrDomainNotFound is an error when no scenario was found.
	ErrDomainNotFound = errors.New("domain not found")
	// ErrScenarioNotFound is an error when no scenario was found.
	ErrScenarioNotFound = errors.New("scenario not found")
)

// 422
var (
	// ErrInvalidBody is an error when provided an invalid request body that can't be parsed.
	ErrInvalidBody = errors.New("can't parse invalid request body")
	// ErrInvalidParams is an error when provided invalid path or query param that can't be parsed.
	ErrInvalidParams = errors.New("can't parse invalid path or query params")
)
