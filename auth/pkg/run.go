package pkg

import (
	"context"
	"github.com/larek-tech/diploma/auth/config"
	server "github.com/larek-tech/diploma/auth/internal/_server"
	"github.com/larek-tech/diploma/auth/internal/auth/controller"
	"github.com/larek-tech/diploma/auth/internal/auth/handler"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/larek-tech/diploma/auth/internal/auth/repo"
	"github.com/larek-tech/diploma/auth/pkg/jwt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
)

const (
	configPath = "./config/config.yaml"
)

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

	provider, err := tracing.NewTraceProvider(exporter, "auth")
	if err != nil {
		return errs.WrapErr(err)
	}
	defer func() {
		if e := provider.Shutdown(ctx); err != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("shutdown tracing provider")
		}
	}()

	otel.SetTracerProvider(provider)
	tracer := otel.Tracer("auth")

	pg, err := postgres.New(&cfg.Postgres, tracer)
	if err != nil {
		return errs.WrapErr(err)
	}
	defer pg.Close()

	authRepo := repo.New(pg)
	jwtProvider := jwt.New(cfg.Jwt)
	authController := controller.New(tracer, authRepo, jwtProvider)
	authHandler := handler.New(tracer, authController)

	srv := server.New(cfg.Server, pg, tracer)
	pb.RegisterAuthServiceServer(srv.GetSrv(), authHandler)
	srv.Start()

	return nil
}
