package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larek-tech/diploma/data/internal/domain/site/service/crawler"

	pageStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/page"

	siteStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/site"

	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_page"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/parse_site"
	"github.com/larek-tech/diploma/data/internal/worker/qaas/result_message"
	"github.com/larek-tech/storage/postgres"
	"go.dataddo.com/pgq"
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
	httpClient := &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			slog.Debug("Not following redirect", "to", req.URL.String())
			return http.ErrUseLastResponse
		},
		Jar: nil,
	}

	pub := qaas.NewPublisher(pgq.NewPublisher(sqlDB))

	siteStore := siteStorage.New(pg)

	pageStore := pageStorage.New(pg)
	pageService := crawler.New(httpClient, siteStore, pageStore, trManager)

	consumer := qaas.NewConsumer(
		parse_page.New(pageStore, pageService, pub),
		parse_site.New(siteStore, pub),
		result_message.New(),
		sqlDB,
	)

	if err = consumer.Run(ctx); err != nil {
		slog.Error("failed to run consumer", "error", err)
		return 1
	}
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

//sourceService "github.com/larek-tech/diploma/data/internal/domain/source/service"
//sourceStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/source"
//sourceStore := sourceStorage.New(pg)
//srcService := sourceService.New(sourceStore, pub, trManager)
