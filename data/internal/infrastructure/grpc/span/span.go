package span

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func GetTraceCtx(ctx context.Context) (context.Context, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	traceIDVal, ok := md["x-trace-id"]
	if !ok {
		slog.Error("trace id not found")
		return nil, nil
	}

	traceIDString := traceIDVal[0]
	traceID, err := trace.TraceIDFromHex(traceIDString)
	if err != nil {
		return ctx, errors.Wrap(err, "failed to parse trace id")
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
	})
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)
	return ctx, nil
}
