package crawler

import (
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	httpClient   httpClient
	siteStore    siteStore
	pageStore    pageStore
	pageJobStore pageJobStore
	trManager    transactionalManager
	tracer       trace.Tracer
}

func New(
	httpClient httpClient,
	siteStorage siteStore,
	pageStorage pageStore,
	pageJobStore pageJobStore,
	trManager transactionalManager,
	tracer trace.Tracer,
) *Service {
	return &Service{
		httpClient:   httpClient,
		siteStore:    siteStorage,
		pageStore:    pageStorage,
		pageJobStore: pageJobStore,
		trManager:    trManager,
		tracer:       tracer,
	}
}
