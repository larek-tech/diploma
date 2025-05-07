package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	documentService "github.com/larek-tech/diploma/data/internal/domain/document/service"
	"github.com/larek-tech/diploma/data/internal/domain/site/service/crawler"
	"github.com/larek-tech/diploma/data/internal/infrastructure/kafka"
	"github.com/larek-tech/diploma/data/internal/infrastructure/ollama"
	chunkStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/chunk"
	documentStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/document"
	pageStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/page"
	questionStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/question"
	siteStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/site"
	"github.com/larek-tech/diploma/data/internal/infrastructure/postgres/sitejob"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/embed_document"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_page"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_site"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_site_status"
	"github.com/larek-tech/storage/postgres"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	var pgCfg postgres.Cfg
	if err := cleanenv.ReadEnv(&pgCfg); err != nil {
		slog.Error("failed to read env", "error", err)
	}
	endpoint, _ := getLLMConfig()
	pg, trManager, err := postgres.New(ctx, pgCfg)
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
	llm, err := ollama.New(endpoint)
	if err != nil {
		slog.Error("failed to create LLM", "error", err)
		return -1
	}
	embedderURL, _ := getEmbedderConfig()
	embedderModel, err := ollama.New(embedderURL)
	if err != nil {
		slog.Error("failed to create embedder", "error", err)
		return -1
	}
	pub := qaas.NewPublisher(sqlDB)
	pub.CreateAllTables([]qaas.Queue{
		qaas.ParseSiteQueue,
		qaas.ParsePageResultQueue,
		qaas.ParsePageQueue,
		qaas.EmbedResultQueue,
		qaas.ParseSiteStatusQueue,
	})

	siteStore := siteStorage.New(pg)
	documentStore := documentStorage.New(pg)
	questionStore := questionStorage.New(pg, trManager)
	chunkStore := chunkStorage.New(pg, trManager)
	pageStore := pageStorage.New(pg)
	siteJobStore := sitejob.New(pg)
	pageService := crawler.New(httpClient, siteStore, pageStore, siteJobStore, trManager)
	embeddingService := documentService.New(documentStore, chunkStore, questionStore, embedderModel, llm, trManager)
	consumer := qaas.NewConsumer(sqlDB)

	slog.Info("Starting consumer")
	wg := &sync.WaitGroup{}
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
		err = consumer.Run(ctx, qaas.ParsePageQueue, parse_page.New(pageStore, pageService, pub))
		if err != nil {
			slog.Error("failed to run consumer", "error", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = consumer.Run(ctx, qaas.ParsePageResultQueue, embed_document.New(embeddingService, pageStore, siteStore))
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

func getLLMConfig() (string, string) {
	host := os.Getenv("OLLAMA_LLM_ENDPOINT")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_LLM_MODEL")
	if model == "" {
		model = "llama3:latest"
	}
	return host, model
}

func getEmbedderConfig() (string, string) {
	host := os.Getenv("OLLAMA_EMBEDDER_ENDPOINT")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_EMBEDDER_MODEL")
	if model == "" {
		model = "bge-m3:latest"
	}
	return host, model
}
