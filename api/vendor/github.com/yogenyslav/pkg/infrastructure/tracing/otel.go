// Package tracing provides a wrapper around OpenTelemetry tracing.
package tracing

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	// ErrNewExporter is an error when a new OTLP exporter can't be created.
	ErrNewExporter = errors.New("failed to create new otlp exporter")
	// ErrNewProvider is an error when a new trace provider can't be created.
	ErrNewProvider = errors.New("new trace provider can't be created")
)

// SetupOTel setups OpenTelemetry tracing with the given configuration.
func SetupOTel(cfg *Config, name string) error {
	exporter, err := NewExporter(context.Background(), cfg.URL())
	if err != nil {
		return err
	}

	provider, err := NewTraceProvider(exporter, name)
	if err != nil {
		return err
	}
	otel.SetTracerProvider(provider)
	return nil
}

// NewExporter creates a new OTLP trace exporter.
func NewExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure(), otlptracehttp.WithEndpoint(endpoint))
	if err != nil {
		return nil, errors.Join(ErrNewExporter, err)
	}
	return exporter, nil
}

// NewTraceProvider creates a new OTel trace provider.
func NewTraceProvider(exp sdktrace.SpanExporter, name string) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
		),
	)
	if err != nil {
		return nil, errors.Join(ErrNewProvider, err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	), nil
}
