package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larek-tech/diploma/data/internal/domain/document/service"
	chunkStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/chunk"
	documentStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/document"
	questionStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/question"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/storage/postgres"
	"go.dataddo.com/pgq/x/schema"
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
	endpoint, model := os.Getenv("OLLAMA_ENDPOINT"), os.Getenv("OLLAMA_MODEL")

	documentStore := documentStorage.New(pg)
	questionStore := questionStorage.New(pg, trManager)
	chunkStore := chunkStorage.New(pg, trManager)

	embeddingService := service.New(documentStore, chunkStore, questionStore, llm, trManager)

	return 0
}

func getSqlCon(db *postgres.DB) *sql.DB {
	pool := db.GetPool()
	sqlCon := stdlib.OpenDBFromPool(pool)
	create := schema.GenerateCreateTableQuery(qaas.QueueName)

	_, err := sqlCon.Exec(create)
	if err != nil {
		slog.Error("Failed to create table", "error", err)
		return nil
	}
	return sqlCon
}
