package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/shared"
	"github.com/yogenyslav/pkg/response"
)

var (
	errMap = map[error]response.ErrorResponse{
		// 400
		shared.ErrWsProtocolRequired: {
			Msg:    "upgrade to ws protocol required",
			Status: fiber.StatusBadRequest,
		},
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
		shared.ErrCreateDomain: {
			Msg:    "failed creating domain",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrGetDomain: {
			Msg:    "failed getting domain",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrUpdateDomain: {
			Msg:    "failed updating domain",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrDeleteDomain: {
			Msg:    "failed deleting domain",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrListDomains: {
			Msg:    "failed listing available domains",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrCreateScenario: {
			Msg:    "failed creating scenario",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrGetScenario: {
			Msg:    "failed getting scenario",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrUpdateScenario: {
			Msg:    "failed updating scenario",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrDeleteScenario: {
			Msg:    "failed deleting scenario",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrListScenarios: {
			Msg:    "failed listing available scenarios",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrCreateChat: {
			Msg:    "failed creating chat",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrGetChat: {
			Msg:    "failed getting chat",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrUpdateChat: {
			Msg:    "failed updating chat",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrDeleteChat: {
			Msg:    "failed deleting chat",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrListChats: {
			Msg:    "failed listing available chats",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrCancelQuery: {
			Msg:    "failed to cancel query processing",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrCreateUser: {
			Msg:    "failed creating user",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrGetUser: {
			Msg:    "failed getting user",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrDeleteUser: {
			Msg:    "failed deleting user",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrListUsers: {
			Msg:    "failed listing available users",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrCreateRole: {
			Msg:    "failed creating role",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrGetRole: {
			Msg:    "failed getting role",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrDeleteRole: {
			Msg:    "failed deleting role",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrListRoles: {
			Msg:    "failed listing available roles",
			Status: fiber.StatusBadRequest,
		},
		shared.ErrUpdateRoleForUser: {
			Msg:    "failed to set or remove user role",
			Status: fiber.StatusBadRequest,
		},
		// 401
		shared.ErrUnauthorized: {
			Msg:    "unauthorized",
			Status: fiber.StatusUnauthorized,
		},
		// 403
		shared.ErrForbidden: {
			Msg:    "user has no access to the resource",
			Status: fiber.StatusForbidden,
		},
		// 404
		shared.ErrSourceNotFound: {
			Msg:    "source not found",
			Status: fiber.StatusNotFound,
		},
		shared.ErrDomainNotFound: {
			Msg:    "domain not found",
			Status: fiber.StatusNotFound,
		},
		shared.ErrScenarioNotFound: {
			Msg:    "scenario not found",
			Status: fiber.StatusNotFound,
		},
		shared.ErrChatNotFound: {
			Msg:    "chat not found",
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
