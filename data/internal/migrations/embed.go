package migrations

import (
	"embed"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type db interface {
	GetPool() *pgxpool.Pool
}

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(db db) func() bool {
	readyCh := make(chan bool)
	go func() {
		goose.SetBaseFS(embedMigrations)
		if err := goose.SetDialect("postgres"); err != nil {
			slog.Error("migration dialect error", "error", err)
			readyCh <- false
			return
		}

		pool := db.GetPool()
		sqlDB := stdlib.OpenDBFromPool(pool)

		if err := goose.Up(sqlDB, "."); err != nil {
			slog.Error("migration failed", "error", err)
			readyCh <- false
			return
		}
		slog.Info("successfully migrated to the latest migration")
		readyCh <- true
	}()

	return func() bool {
		select {
		case result := <-readyCh:
			return result
		default:
			return false
		}
	}
}
