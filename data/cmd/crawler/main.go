package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	sourceService "github.com/larek-tech/diploma/data/internal/domain/source/service"
	chunkStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/chunk"
	sourceStorage "github.com/larek-tech/diploma/data/internal/infrastructure/postgres/source"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/storage/postgres"
	"github.com/tmc/langchaingo/llms/ollama"
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
	pub := qaas.NewPublisher(pgq.NewPublisher(sqlDB))

	sourceStore := sourceStorage.New(pg)
	srcService := sourceService.New(sourceStore, pub, trManager)

	chunkStore := chunkStorage.New(pg, trManager)
	embedder, err := ollama.New(
		ollama.WithServerURL("http://127.0.0.1:11434"),
		ollama.WithModel("bge-m3"),
	)

	http.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload source.DataMessage
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		src, err := srcService.CreateSource(ctx, payload)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(src)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}))

	http.Handle("/q", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload SearchQuery
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		if len(payload.SourceIDs) == 0 {
			http.Error(w, "Source IDs are required", http.StatusBadRequest)
			return
		}
		if payload.TopK == 0 {
			payload.TopK = 5
		}
		if payload.Threshold == 0 {
			payload.Threshold = 0.1
		}
		embededQuery, err := embedder.CreateEmbedding(ctx, []string{payload.Query})
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		res, err := chunkStore.Search(ctx, embededQuery[0], payload.SourceIDs, int(payload.TopK))
		if err != nil {
			http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(res) == 0 {
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}
		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, fmt.Sprintf("marshaling err: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))

	slog.Info("Starting server on :8080")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("Failed to start server", "error", err)
		return 1
	}
	return 0
}

// string query = 1;
// repeated string sourceIds = 2;
// uint64 topK = 3;
// float threshold = 4;
// bool useQuestions = 5; // hypothetical questions
type SearchQuery struct {
	Query     string   `json:"query"`
	SourceIDs []string `json:"sourceIds"`
	TopK      uint     `json:"topK"`
	Threshold float32  `json:"threshold"`
	UseQ      bool     `json:"useQuestions"`
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
