package metric

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	sourcesCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sources_created",
			Help: "Number of sources created",
		},
		[]string{"source_type", "source_id", "err"},
	)
	documentsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "documents_processed",
			Help: "Number of documents processed",
		},
		[]string{"source_type", "source_id", "err"},
	)
	documentsParsed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "documents_parsed",
			Help: "Number of documents parsed",
		},
		[]string{"object_id", "source_type", "source_id", "err"},
	)
	chunksCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chunks_created",
			Help: "Number of chunks created",
		},
		[]string{"document_id", "source_type", "source_id", "err"},
	)
	searchQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "search_queries",
			Help: "Number of search queries",
		},
		[]string{"source_id", "source_type", "query", "err"},
	)
)

func InitializeMetrics() {
	prometheus.MustRegister(sourcesCreated)
	prometheus.MustRegister(documentsProcessed)
	prometheus.MustRegister(documentsParsed)
	prometheus.MustRegister(chunksCreated)
	prometheus.MustRegister(searchQueries)
}

func IncrementSourcesCreated(sourceType, sourceID string, err error) {
	errStr := errToString(err)
	sourcesCreated.WithLabelValues(sourceType, sourceID, errStr).Inc()
}

func IncrementDocumentsProcessed(sourceType, sourceID string, err error) {
	errStr := errToString(err)
	documentsProcessed.WithLabelValues(sourceType, sourceID, errStr).Inc()
}

func IncrementDocumentsParsed(objectID, sourceType, sourceID string, err error) {
	errStr := errToString(err)
	documentsParsed.WithLabelValues(objectID, sourceType, sourceID, errStr).Inc()
}

func IncrementChunksCreated(documentID, sourceID, sourceType string, err error, cnt ...int) {
	errStr := errToString(err)
	if len(cnt) > 0 {
		chunksCreated.WithLabelValues(documentID, sourceID, sourceType, errStr).Add(float64(cnt[0]))
		return
	}
	chunksCreated.WithLabelValues(documentID, sourceID, sourceType, errStr).Inc()
}

func IncrementSearchQueries(sourceID, sourceType, query string, err error) {
	errStr := errToString(err)
	searchQueries.WithLabelValues(sourceID, sourceType, query, errStr).Inc()
}

func errToString(err error) string {
	if err != nil {
		return "true"
	}
	return "false"
}

func RunPrometheusServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting Prometheus server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("Failed to start Prometheus server", "err", err)
	}
}
