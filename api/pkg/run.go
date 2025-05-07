package pkg

import (
	"context"

	"github.com/larek-tech/diploma/api/config"
	server "github.com/larek-tech/diploma/api/internal/_server"
	"github.com/larek-tech/diploma/api/internal/api"
	"github.com/larek-tech/diploma/api/internal/auth"
	"github.com/larek-tech/diploma/api/internal/auth/handler"
	"github.com/larek-tech/diploma/api/internal/auth/middleware"
	"github.com/larek-tech/diploma/api/internal/auth/pb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	grpcclient "github.com/yogenyslav/pkg/grpc_client"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	configPath = "./config/config.yaml"
)

// Run setup application and run it.
func Run() error {
	cfg, err := config.New(configPath)
	if err != nil {
		return errs.WrapErr(err)
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errs.WrapErr(err)
	}
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	exporter, err := tracing.NewExporter(ctx, cfg.Jaeger.URL())
	if err != nil {
		return errs.WrapErr(err)
	}
	defer func() {
		if e := exporter.Shutdown(ctx); e != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("shutdown tracing exporter")
		}
	}()

	provider, err := tracing.NewTraceProvider(exporter, "api")
	if err != nil {
		return errs.WrapErr(err)
	}
	defer func() {
		if e := provider.Shutdown(ctx); err != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("shutdown tracing provider")
		}
	}()

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := otel.Tracer("api")

	pg, err := postgres.New(&cfg.Postgres, tracer)
	if err != nil {
		return errs.WrapErr(err)
	}
	defer pg.Close()

	// Connect to auth service via gRPC
	authConn, err := grpcclient.NewGrpcClient(
		&cfg.AuthService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return errs.WrapErr(err, "connect to auth service")
	}
	defer authConn.Close()
	authService := pb.NewAuthServiceClient(authConn.Conn())

	srv := server.New(cfg.Server)

	// Auth routes
	authRouter := srv.GetSrv().Group("/auth/v1")
	authHandler := handler.New(authService, tracer)
	auth.SetupRoutes(authRouter, authHandler)

	// Connect to domain service via gRPC
	domainConn, err := grpcclient.NewGrpcClient(
		&cfg.DomainService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return errs.WrapErr(err, "connect to domain service")
	}
	defer domainConn.Close()

	// Api routes with JWT middleware
	apiRouter := srv.GetSrv().Group("/api/v1")
	apiRouter.Use(middleware.Jwt(authService))
	api.SetupRoutes(apiRouter, domainConn.Conn())

	srv.Start()

	return nil
}
