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
	sc "github.com/larek-tech/diploma/domain/internal/domain/source/controller"
	sh "github.com/larek-tech/diploma/domain/internal/domain/source/handler"
	sr "github.com/larek-tech/diploma/domain/internal/domain/source/repo"
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

	srv.Start()

	return nil
}
