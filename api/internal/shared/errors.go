package shared

import (
	"errors"
)

// 400
var (
	// ErrWsProtocolRequired is an error when required to upgrade to ws protocol.
	ErrWsProtocolRequired = errors.New("upgrade to ws protocol required")

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

	// ErrCreateChat is an error when failed to create chat.
	ErrCreateChat = errors.New("failed to create chat")
	// ErrGetChat is an error when failed to get chat.
	ErrGetChat = errors.New("failed to get chat")
	// ErrUpdateChat is an error when failed to update chat.
	ErrUpdateChat = errors.New("failed to update chat")
	// ErrDeleteChat is an error when failed to delete chat.
	ErrDeleteChat = errors.New("failed to delete chat")
	// ErrListChats is an error when failed to list chats.
	ErrListChats = errors.New("failed to list chats")
	// ErrCancelQuery is an error when failed to cancel processing query.
	ErrCancelQuery = errors.New("failed to cancel query")

	// ErrCreateUser is an error when failed to create user.
	ErrCreateUser = errors.New("failed to create user")
	// ErrGetUser is an error when failed to get user.
	ErrGetUser = errors.New("failed to get user")
	// ErrDeleteUser is an error when failed to delete user.
	ErrDeleteUser = errors.New("failed to delete user")
	// ErrListUsers is an error when failed to list users.
	ErrListUsers = errors.New("failed to list users")

	// ErrCreateRole is an error when failed to create role.
	ErrCreateRole = errors.New("failed to create role")
	// ErrGetRole is an error when failed to get role.
	ErrGetRole = errors.New("failed to get role")
	// ErrDeleteRole is an error when failed to delete role.
	ErrDeleteRole = errors.New("failed to delete role")
	// ErrListRoles is an error when failed to list roles.
	ErrListRoles = errors.New("failed to list roles")
	// ErrUpdateRoleForUser is an error when failed to set/remove role for user.
	ErrUpdateRoleForUser = errors.New("failed to update role for user")
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
	// ErrChatNotFound is an error when no chat was found.
	ErrChatNotFound = errors.New("chat not found")
)

// 422
var (
	// ErrInvalidBody is an error when provided an invalid request body that can't be parsed.
	ErrInvalidBody = errors.New("can't parse invalid request body")
	// ErrInvalidParams is an error when provided invalid path or query param that can't be parsed.
	ErrInvalidParams = errors.New("can't parse invalid path or query params")
)
