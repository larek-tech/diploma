package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/migrations"
	"github.com/larek-tech/storage/postgres"
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
	pub := qaas.NewPublisher(stdlib.OpenDBFromPool(db.GetPool()))
	pub.CreateAllTables([]qaas.Queue{
		qaas.ParseSiteQueue,
		qaas.ParsePageResultQueue,
		qaas.ParsePageQueue,
		qaas.EmbedResultQueue,
	})

	slog.Info("migration finished")
	return 0
}
