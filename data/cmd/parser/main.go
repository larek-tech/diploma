package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	documentService "github.com/larek-tech/diploma/data/internal/domain/document/service"
	questionService "github.com/larek-tech/diploma/data/internal/domain/question/service"
	"github.com/larek-tech/diploma/data/internal/domain/site/service/crawler"
	"github.com/larek-tech/diploma/data/internal/infrastructure/kafka"
	"github.com/larek-tech/diploma/data/internal/infrastructure/ocr"
	"github.com/larek-tech/diploma/data/internal/infrastructure/ollama"
	"github.com/larek-tech/diploma/data/internal/infrastructure/s3"
	chunkStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/chunk"
	documentStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/document"
	fileStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/file"
	pageStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/page"
	questionStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/question"
	siteStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/storage/sitejob"
	"github.com/larek-tech/diploma/data/pkg/metric"
	"github.com/otiai10/gosseract"
	"github.com/yogenyslav/pkg/infrastructure/tracing"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/embed_document"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_page"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_site"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_site_status"
	"github.com/larek-tech/storage/postgres"
)

const (
	serviceName = "parser"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	exporter, err := tracing.NewExporter(ctx, getTracingEndpoint())
	if err != nil {
		slog.Error(err.Error())
		return -1
	}
	defer func() {
		if e := exporter.Shutdown(ctx); e != nil {
			slog.Error("shutdown tracing exporter", "error", e)
		}
	}()
	provider, err := tracing.NewTraceProvider(exporter, serviceName)
	if err != nil {
		slog.Error(err.Error())
		return -1
	}
	tracer := provider.Tracer(serviceName)

	var pgCfg postgres.Cfg
	if err := cleanenv.ReadEnv(&pgCfg); err != nil {
		slog.Error("failed to read env", "error", err)
	}

	pg, trManager, err := postgres.New(ctx, pgCfg,
		postgres.WithTelemetry(true),
		postgres.WithTracer(tracer),
	)
	if err != nil {
		slog.Error("failed to create postgres", "error", err)
		return 1
	}
	defer pg.Close()
	sqlDB := getSqlCon(pg)
	if sqlDB == nil {
		slog.Error("Failed to get SQL connection")
		return 1
	}
	var kafkaCfg kafka.Config
	if err := cleanenv.ReadEnv(&kafkaCfg); err != nil {
		slog.Error("failed to read env", "error", err)
		return -1
	}
	kafkaProducer, err := kafka.NewProducer(kafkaCfg)
	if err != nil {
		slog.Error("failed to create kafka producer", "error", err)
		return -1
	}

	httpClient := &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			slog.Debug("Not following redirect", "to", req.URL.String())
			return http.ErrUseLastResponse
		},
		Jar: nil,
	}
	endpoint, llmModel, contextSize := getLLMConfig()
	llm, err := ollama.New(endpoint, &ollama.Config{
		LLMModel:       llmModel,
		LLMContextSize: contextSize,
	})
	if err != nil {
		slog.Error("failed to create LLM", "error", err)
		return -1
	}
	embedderURL, embedderModel, embeddingsSize := getEmbedderConfig()
	embedderService, err := ollama.New(embedderURL, &ollama.Config{
		EmbeddingSize:   embeddingsSize,
		EmbeddingsModel: embedderModel,
	})
	if err != nil {
		slog.Error("failed to create embedder", "error", err)
		return -1
	}
	pub := qaas.NewPublisher(sqlDB)
	err = pub.CreateAllTables([]qaas.Queue{
		qaas.ParseSiteQueue,
		qaas.ParsePageResultQueue,
		qaas.ParsePageQueue,
		qaas.EmbedResultQueue,
		qaas.ParseSiteStatusQueue,
	})
	if err != nil {
		slog.Error("failed to create tables", "error", err)
		return -1
	}

	objectStorage, err := s3.New(getS3Credentials())
	if err != nil {
		slog.Error("failed to create s3 client", "error", err)
		return -1
	}
	err = objectStorage.CreateBuckets(ctx,
		pageStorage.PageBucketName,
		fileStorage.FileBucketName,
	)
	if err != nil {
		slog.Error("failed to create s3 bucket", "error", err)
		return -1
	}

	tesseract := gosseract.NewClient()
	tesseract.Languages = []string{"rus", "eng"}
	defer tesseract.Close()

	ocr := ocr.New(tesseract)
	fileStorage := fileStorage.New(pg, objectStorage)
	siteStore := siteStorage.New(pg)
	documentStore := documentStorage.New(pg)
	questionStore := questionStorage.New(pg, trManager)
	chunkStore := chunkStorage.New(pg, trManager)
	pageStore := pageStorage.New(pg, objectStorage)
	siteJobStore := sitejob.New(pg)
	pageService := crawler.New(httpClient, siteStore, pageStore, siteJobStore, trManager, tracer)
	questionSrv := questionService.New(llm, embedderService)
	embeddingService := documentService.New(documentStore, chunkStore, questionStore, questionSrv, embedderService, ocr, trManager, tracer)
	consumer := qaas.NewConsumer(sqlDB)

	slog.Info("Starting consumer")
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		metric.RunPrometheusServer("9091")
	}()

	wg.Add(1)
	// site parser
	go func() {
		defer wg.Done()
		err = consumer.Run(ctx, qaas.ParseSiteQueue, parse_site.New(siteStore, pub))
		if err != nil {
			slog.Error("failed to run consumer", "error", err)
		}
	}()
	// site page parser
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = consumer.Run(ctx, qaas.ParsePageQueue, parse_page.New(pageStore, pageService, pub, tracer))
		if err != nil {
			slog.Error("failed to run consumer", "error", err)
		}
	}()
	wg.Add(1)

	go func() {
		defer wg.Done()
		err = consumer.Run(ctx, qaas.ParsePageResultQueue, embed_document.New(embeddingService, pageStore, siteStore, fileStorage))
		if err != nil {
			slog.Error("failed to run consumer", "error", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = consumer.Run(ctx, qaas.ParseFileQueue, embed_document.New(embeddingService, pageStore, siteStore, fileStorage))
		if err != nil {
			slog.Error("failed to run consumer", "error", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = consumer.Run(ctx, qaas.ParseSiteStatusQueue, parse_site_status.New(pub, siteJobStore, kafkaProducer))
		if err != nil {
			slog.Error("failed to run consumer", "error", err)
		}
	}()

	wg.Wait()
	return 0
}

func getSqlCon(db *postgres.DB) *sql.DB {
	pool := db.GetPool()
	sqlCon := stdlib.OpenDBFromPool(pool)
	return sqlCon
}

func getEmbedderConfig() (string, string, int) {
	host := os.Getenv("OLLAMA_EMBEDDER_ENDPOINT")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_EMBEDDER_MODEL")
	if model == "" {
		model = "bge-m3:latest"
	}
	embeddingSize := os.Getenv("OLLAMA_EMBEDDER_SIZE")
	size, err := strconv.Atoi(embeddingSize)
	if err != nil {
		size = 514
	}

	return host, model, size
}

func getLLMConfig() (string, string, int) {
	host := os.Getenv("OLLAMA_LLM_ENDPOINT")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_LLM_MODEL")
	if model == "" {
		model = "llama3:latest"
	}
	contextSize := os.Getenv("OLLAMA_LLM_CONTEXT_SIZE")
	size, err := strconv.Atoi(contextSize)
	if err != nil {
		size = 32000
	}

	return host, model, size
}

func getS3Credentials() s3.Credentials {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	return s3.NewCredentials(endpoint, os.Getenv("S3_ACCESS_KEY_ID"), os.Getenv("S3_SECRET_ACCESS_KEY"), true)
}

func getTracingEndpoint() string {
	tracingEndpoint := os.Getenv("TRACING_ENDPOINT")
	if tracingEndpoint == "" {
		tracingEndpoint = "localhost:4318"
	}
	return tracingEndpoint
}
