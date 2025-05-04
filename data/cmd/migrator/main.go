package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/migrations"
	"github.com/larek-tech/storage/postgres"
	"go.dataddo.com/pgq/x/schema"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	var cfg postgres.Cfg

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("failed to read env", "error", err)
		return 1
	}
	db, _, err := postgres.New(ctx, cfg)
	if err != nil {
		slog.Error("failed to connect to db", "error", err)
		return 1
	}
	defer db.Close()

	ready := migrations.Migrate(db)
	for !ready() {
		slog.Info("waiting for migration to finish")
		time.Sleep(time.Second * 1)
	}
	create := schema.GenerateCreateTableQuery(qaas.QueueName)
	if err = db.Exec(ctx, create); err != nil {
		slog.Error("failed to create pgq", "error", err)
		return 1
	}

	slog.Info("migration finished")
	return 0
}
