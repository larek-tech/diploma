package api

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/larek-tech/diploma/api/internal/api/chat"
	ch "github.com/larek-tech/diploma/api/internal/api/chat/handler"
	"github.com/larek-tech/diploma/api/internal/api/domain"
	dh "github.com/larek-tech/diploma/api/internal/api/domain/handler"
	"github.com/larek-tech/diploma/api/internal/api/role"
	rh "github.com/larek-tech/diploma/api/internal/api/role/handler"
	"github.com/larek-tech/diploma/api/internal/api/scenario"
	sch "github.com/larek-tech/diploma/api/internal/api/scenario/handler"
	"github.com/larek-tech/diploma/api/internal/api/source"
	sh "github.com/larek-tech/diploma/api/internal/api/source/handler"
	"github.com/larek-tech/diploma/api/internal/api/user"
	uh "github.com/larek-tech/diploma/api/internal/api/user/handler"
	authpb "github.com/larek-tech/diploma/api/internal/auth/pb"
	chatpb "github.com/larek-tech/diploma/api/internal/chat/pb"
	domainpb "github.com/larek-tech/diploma/api/internal/domain/pb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// SetupRoutes maps api routes.
func SetupRoutes(
	api fiber.Router,
	domainConn, chatConn, authConn, mlConn *grpc.ClientConn,
	tracer trace.Tracer,
	wsConfig websocket.Config,
) {
	sourceRouter := api.Group("/source")
	sourceHandler := sh.New(domainpb.NewSourceServiceClient(domainConn))
	source.SetupRoutes(sourceRouter, sourceHandler)

	domainRouter := api.Group("/domain")
	domainHandler := dh.New(
		domainpb.NewDomainServiceClient(domainConn),
		domainpb.NewScenarioServiceClient(domainConn),
		domainpb.NewSourceServiceClient(domainConn),
		domainpb.NewMLServiceClient(mlConn),
	)
	domain.SetupRoutes(domainRouter, domainHandler)

	scenarioRouter := api.Group("/scenario")
	scenarioHandler := sch.New(domainpb.NewScenarioServiceClient(domainConn))
	scenario.SetupRoutes(scenarioRouter, scenarioHandler)

	chatRouter := api.Group("/chat")
	chatHandler := ch.New(
		chatpb.NewChatServiceClient(chatConn),
		authpb.NewAuthServiceClient(authConn),
		domainpb.NewMLServiceClient(mlConn),
		domainpb.NewScenarioServiceClient(domainConn),
		domainpb.NewDomainServiceClient(domainConn),
		domainpb.NewSourceServiceClient(domainConn),
		tracer,
	)
	chat.SetupRoutes(chatRouter, chatHandler, wsConfig)

	userRouter := api.Group("/user")
	userHandler := uh.New(domainpb.NewUserServiceClient(domainConn))
	user.SetupRoutes(userRouter, userHandler)

	roleRouter := api.Group("/role")
	roleHandler := rh.New(domainpb.NewRoleServiceClient(domainConn))
	role.SetupRoutes(roleRouter, roleHandler)
}
