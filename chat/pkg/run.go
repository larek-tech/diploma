package pkg

import (
	"context"

	"github.com/larek-tech/diploma/chat/config"
	server "github.com/larek-tech/diploma/chat/internal/_server"
	"github.com/larek-tech/diploma/chat/internal/chat/controller"
	"github.com/larek-tech/diploma/chat/internal/chat/handler"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"github.com/larek-tech/diploma/chat/internal/chat/repo"
	mlpb "github.com/larek-tech/diploma/chat/internal/domain/pb"
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

	provider, err := tracing.NewTraceProvider(exporter, "chat")
	if err != nil {
		return errs.WrapErr(err)
	}
	defer func() {
		if e := provider.Shutdown(ctx); e != nil {
			log.Warn().Err(errs.WrapErr(e)).Msg("shutdown tracing provider")
		}
	}()

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := otel.Tracer("chat")

	pg, err := postgres.New(&cfg.Postgres, tracer)
	if err != nil {
		return errs.WrapErr(err)
	}
	defer pg.Close()

	mlConn, err := grpcclient.NewGrpcClient(
		&cfg.MLService,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return errs.WrapErr(err, "create ml service client")
	}

	srv := server.New(cfg.Server)

	chatRepo := repo.New(pg)
	chatController := controller.New(chatRepo, tracer, mlpb.NewMLServiceClient(mlConn.Conn()))
	chatHandler := handler.New(chatController, tracer)
	pb.RegisterChatServiceServer(srv.GetSrv(), chatHandler)

	srv.Start()

	return nil
}
