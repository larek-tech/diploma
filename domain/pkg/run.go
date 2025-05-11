package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/larek-tech/diploma/domain/config"
	server "github.com/larek-tech/diploma/domain/internal/_server"
	dc "github.com/larek-tech/diploma/domain/internal/domain/domain/controller"
	dh "github.com/larek-tech/diploma/domain/internal/domain/domain/handler"
	dr "github.com/larek-tech/diploma/domain/internal/domain/domain/repo"
	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	rc "github.com/larek-tech/diploma/domain/internal/domain/role/controller"
	rh "github.com/larek-tech/diploma/domain/internal/domain/role/handler"
	rr "github.com/larek-tech/diploma/domain/internal/domain/role/repo"
	scc "github.com/larek-tech/diploma/domain/internal/domain/scenario/controller"
	sch "github.com/larek-tech/diploma/domain/internal/domain/scenario/handler"
	scr "github.com/larek-tech/diploma/domain/internal/domain/scenario/repo"
	sc "github.com/larek-tech/diploma/domain/internal/domain/source/controller"
	sh "github.com/larek-tech/diploma/domain/internal/domain/source/handler"
	sr "github.com/larek-tech/diploma/domain/internal/domain/source/repo"
	uc "github.com/larek-tech/diploma/domain/internal/domain/user/controller"
	uh "github.com/larek-tech/diploma/domain/internal/domain/user/handler"
	ur "github.com/larek-tech/diploma/domain/internal/domain/user/repo"
	"github.com/larek-tech/diploma/domain/pkg/kafka"
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
		if e := provider.Shutdown(ctx); e != nil {
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

	adminBrokerAddr := fmt.Sprintf("%s:%d", cfg.Kafka.Brokers[0].Host, cfg.Kafka.Brokers[0].Port)
	clusterAdmin, err := sarama.NewClusterAdmin([]string{adminBrokerAddr}, sarama.NewConfig())
	if err != nil {
		return errs.WrapErr(err, "kafka setup")
	}
	for _, topic := range cfg.Kafka.Topics {
		err = clusterAdmin.CreateTopic(topic.Name, &sarama.TopicDetail{
			NumPartitions:     topic.Partitions,
			ReplicationFactor: 1,
		}, false)
		if err != nil && !errors.Is(err, sarama.ErrTopicAlreadyExists) {
			return errs.WrapErr(err, "create kafka topic")
		}
	}

	kafkaProducer, errCh, err := kafka.NewAsyncProducer(&cfg.Kafka, sarama.NewRoundRobinPartitioner, sarama.WaitForAll)
	if err != nil {
		return errs.WrapErr(err, "create kafka producer")
	}
	go func() {
		for e := range errCh {
			log.Warn().Err(errs.WrapErr(e)).Msg("kafka producer error")
		}
	}()
	defer kafkaProducer.Close()

	kafkaConsumer, err := kafka.NewConsumer(&cfg.Kafka)
	if err != nil {
		return errs.WrapErr(err, "create kafka consumer")
	}
	defer kafkaConsumer.SingleConsumer.Close()

	srv := server.New(cfg.Server)

	// Setup source module
	sourceRepo := sr.New(pg)
	sourceController, err := sc.New(ctx, sourceRepo, tracer, kafkaProducer, kafkaConsumer)
	if err != nil {
		return errs.WrapErr(err, "create source controller")
	}
	sourceHandler := sh.New(sourceController, tracer)
	pb.RegisterSourceServiceServer(srv.GetSrv(), sourceHandler)

	// Setup domain module
	domainRepo := dr.New(pg)
	domainController := dc.New(domainRepo, tracer)
	domainHandler := dh.New(domainController, tracer)
	pb.RegisterDomainServiceServer(srv.GetSrv(), domainHandler)

	// Setup scenario module
	scenarioRepo := scr.New(pg)
	scenarioController := scc.New(scenarioRepo, tracer)
	scenarioHandler := sch.New(scenarioController, tracer)
	pb.RegisterScenarioServiceServer(srv.GetSrv(), scenarioHandler)

	// Setup user module
	userRepo := ur.New(pg)
	userController := uc.New(userRepo, tracer, cfg.Server.Encryption)
	userHandler := uh.New(userController)
	pb.RegisterUserServiceServer(srv.GetSrv(), userHandler)

	roleRepo := rr.New(pg)
	roleController := rc.New(roleRepo, tracer)
	roleHandler := rh.New(roleController)
	pb.RegisterRoleServiceServer(srv.GetSrv(), roleHandler)

	srv.Start()

	return nil
}
