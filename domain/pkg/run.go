package pkg

import (
	"context"

	"github.com/larek-tech/diploma/domain/config"
	server "github.com/larek-tech/diploma/domain/internal/_server"
	"github.com/larek-tech/diploma/domain/internal/domain/handler"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	provider, err := tracing.NewTraceProvider(exporter, "domain")
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
	tracer := otel.Tracer("domain")

	pg, err := postgres.New(&cfg.Postgres, tracer)
	if err != nil {
		return errs.WrapErr(err)
	}
	defer pg.Close()

	srv := server.New(cfg.Server)

	domainHandler := handler.New(nil, tracer)
	pb.RegisterDomainServiceServer(srv.GetSrv(), domainHandler)

	srv.Start()

	return nil
}
