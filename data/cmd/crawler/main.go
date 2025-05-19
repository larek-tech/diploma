package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/larek-tech/diploma/data/internal/data/pb"
	sitemap "github.com/larek-tech/diploma/data/internal/domain/sitemap/service"
	"github.com/larek-tech/diploma/data/internal/domain/source"
	sourceService "github.com/larek-tech/diploma/data/internal/domain/source/service"
	"github.com/larek-tech/diploma/data/internal/grpc/get_documents"
	"github.com/larek-tech/diploma/data/internal/grpc/vector_search"
	"github.com/larek-tech/diploma/data/internal/infrastructure/grpc/server"
	"github.com/larek-tech/diploma/data/internal/infrastructure/kafka"
	"github.com/larek-tech/diploma/data/internal/infrastructure/ollama"
	"github.com/larek-tech/diploma/data/internal/infrastructure/qaas"
	"github.com/larek-tech/diploma/data/internal/infrastructure/s3"
	chunkStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/chunk"
	documentStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/document"
	fileStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/file"
	pageStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/page"
	sourceStorage "github.com/larek-tech/diploma/data/internal/infrastructure/storage/source"
	"github.com/larek-tech/diploma/data/internal/worker/kafka/create_source"
	"github.com/larek-tech/diploma/data/pkg/metric"
	"github.com/larek-tech/storage/postgres"
	"google.golang.org/grpc/reflection"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	slog.Info("Starting server")
	kafkaCfg, err := getKafkaConfig()
	slog.Info("kafka config", "cfg", *kafkaCfg)
	if err != nil {
		slog.Error(err.Error())
		return -1
	}

	pgCfg, err := getPGConfig()
	if err != nil {
		slog.Error(err.Error())
		return -1
	}
	pg, trManager, err := postgres.New(ctx, *pgCfg)
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

	pub := qaas.NewPublisher(sqlDB)
	pub.CreateAllTables([]qaas.Queue{
		qaas.ParseFileQueue,
		qaas.ParseFileResult,
		qaas.ParseSiteQueue,
		qaas.ParsePageResultQueue,
		qaas.ParsePageQueue,
		qaas.ParseS3Queue,
		qaas.ParseS3ResultQueue,
		qaas.EmbedResultQueue,
	})

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

	fileStorage := fileStorage.New(pg, objectStorage)
	sourceStore := sourceStorage.New(pg)
	srcService := sourceService.New(sourceStore, fileStorage, sitemap.New(), pub, trManager)
	documentStore := documentStorage.New(pg)
	chunkStore := chunkStorage.New(pg, trManager)
	host, _ := getOllamaConfig()
	embedder, err := ollama.New(host)
	if err != nil {
		slog.Error("failed to create ollama client", "error", err)
		return 1
	}
	kafkaProducer, err := kafka.NewProducer(*kafkaCfg)
	if err != nil {
		slog.Error("failed to create kafka producer")
		return 1
	}
	kafkaHandlers := map[string]kafka.HandlerFunc{
		"source": create_source.New(srcService, kafkaProducer).Handle,
	}
	kafkaConsumer, err := kafka.NewConsumer(*kafkaCfg, "crawler", kafkaHandlers)
	if err != nil {
		// for now kafka can be disabled and we can accept messages from http and grpc
		slog.Error("failed to create kafka consumer: %w", "err", err)
	}

	http.Handle("/test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		slog.Info("Received test request")
		if r.Method != http.MethodPost {
			logError(w, "Method not allowed", nil, http.StatusMethodNotAllowed)
			return
		}
		var payload source.DataMessage
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			logError(w, "Bad request", err, http.StatusBadRequest)
			return
		}
		src, err := srcService.CreateSource(ctx, payload)
		if err != nil {
			logError(w, "Internal server error", err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp, err := json.Marshal(src)
		if err != nil {
			logError(w, "Internal server error", err, http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}))

	http.Handle("/q", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			slog.Error("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload SearchQuery
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			slog.Error("Failed to decode request body", "error", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		slog.Info("Received query", "query", payload.Query)
		if len(payload.SourceIDs) == 0 {
			slog.Error("Source IDs are required")
			http.Error(w, "Source IDs are required", http.StatusBadRequest)
			return
		}
		if payload.TopK == 0 {
			payload.TopK = 5
		}
		if payload.Threshold == 0 {
			payload.Threshold = 0.1
		}
		embedding, err := embedder.CreateEmbedding(ctx, []string{payload.Query})
		if err != nil {
			slog.Error("Failed to create embedding", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		res, err := chunkStore.Search(ctx, embedding[0], payload.SourceIDs, payload.Threshold, int(payload.TopK))
		if err != nil {
			slog.Error("Failed to search chunks", "error", err)
			http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
			return
		}
		if len(res) == 0 {
			slog.Error("No results found")
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}
		resp, err := json.Marshal(res)
		if err != nil {
			slog.Error("Failed to marshal response", "error", err)
			http.Error(w, fmt.Sprintf("marshaling err: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}))
	wg := sync.WaitGroup{}
	srv := server.New()
	pb.RegisterDataServiceServer(
		srv.GetSrv(),
		server.NewHandlers(
			vector_search.New(chunkStore, embedder),
			get_documents.New(documentStore),
		),
	)
	reflection.Register(srv.GetSrv())
	wg.Add(1)
	go func() {
		defer wg.Done()
		metric.RunPrometheusServer("9090")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("gRPC server started")
		grpcErr := srv.Run()
		if grpcErr != nil {
			slog.Error("Failed to start gRPC server", "error", grpcErr)
			return
		}

	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if kafkaConsumer != nil {
			slog.Info("Starting kafka consumer")
			kafkaErr := kafkaConsumer.Run(ctx)
			if kafkaErr != nil {
				slog.Error("failed to run kafka: %w", "err", kafkaErr)
			}
		} else {
			slog.Info("Kafka consumer is disabled")
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("Starting server on :8080")
		if err = http.ListenAndServe(":8080", nil); err != nil {
			slog.Error("Failed to start server", "error", err)
		}
	}()
	wg.Wait()
	return 0
}

func getS3Credentials() s3.Credentials {
	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	return s3.NewCredentials(endpoint, os.Getenv("S3_ACCESS_KEY_ID"), os.Getenv("S3_SECRET_ACCESS_KEY"), true)
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
	return sqlCon
}

func getOllamaConfig() (string, string) {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "bge-m3:latest"
	}
	return host, model
}

func getPGConfig() (*postgres.Cfg, error) {
	var pgCfg postgres.Cfg
	if err := cleanenv.ReadEnv(&pgCfg); err != nil {
		return nil, fmt.Errorf("failed to read postgresql config: %w", err)

	}
	return &pgCfg, nil
}

func getKafkaConfig() (*kafka.Config, error) {
	var cfg kafka.Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read kafka config: %w", err)
	}
	return &cfg, nil
}

func logError(w http.ResponseWriter, msg string, err error, code int) {
	slog.Error(msg, "error", err)
	http.Error(w, msg, code)
}
