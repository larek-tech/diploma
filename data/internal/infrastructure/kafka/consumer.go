package kafka

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama" // Changed import
	"github.com/samber/lo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const handlerTimeout = 300 * time.Second

// HandlerFunc signature updated for sarama.ConsumerMessage
type HandlerFunc func(context.Context, *sarama.ConsumerMessage) error

// Consumer struct updated for Sarama
type Consumer struct {
	cfg           Config
	consumerGroup sarama.ConsumerGroup
	handlers      map[string]HandlerFunc // map of topic to handler
	groupID       string
	tracer        trace.Tracer
}

// NewConsumer updated for Sarama
// Added groupID parameter as it's essential for Sarama consumer groups.
func NewConsumer(cfg Config, groupID string, handlers map[string]HandlerFunc, tracer trace.Tracer) (*Consumer, error) {
	if groupID == "" {
		return nil, fmt.Errorf("kafka consumer groupID cannot be empty")
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("no handlers provided")
	}

	saramaCfg := sarama.NewConfig()

	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaCfg.Consumer.Return.Errors = true // To receive errors from the consumer group error channel

	consumerGroup, err := sarama.NewConsumerGroup([]string{cfg.BootstrapServers}, groupID, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init sarama consumer group: %w", err)
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		handlers:      handlers,
		cfg:           cfg,
		groupID:       groupID,
		tracer:        tracer,
	}, nil
}

// saramaGroupHandler implements sarama.ConsumerGroupHandler
type saramaGroupHandler struct {
	handlers       map[string]HandlerFunc
	handlerTimeout time.Duration
	tracer         trace.Tracer
}

func (h *saramaGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	slog.Info("Sarama consumer group handler setup", "memberID", session.MemberID(), "claims", session.Claims())
	return nil
}

func (h *saramaGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	slog.Info("Sarama consumer group handler cleanup", "memberID", session.MemberID())
	return nil
}

func getTraceCtx(ctx context.Context, headers []*sarama.RecordHeader) (context.Context, error) {
	var traceIDString string
	for _, h := range headers {
		if string(h.Key) == "x-trace-id" {
			traceIDString = string(h.Value)
			break
		}
	}
	if traceIDString == "" {
		return nil, fmt.Errorf("failed to get trace id")
	}
	traceID, err := trace.TraceIDFromHex(traceIDString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trace id: %w", err)
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
	})
	return trace.ContextWithSpanContext(ctx, spanCtx), nil
}

func (h *saramaGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			ctx := context.Background()
			if !ok {
				slog.Info("Message channel closed for claim", "topic", claim.Topic(), "partition", claim.Partition())
				return nil // Exit when messages channel is closed
			}
			if message == nil {
				// This case should ideally not happen if ok is true, but good for robustness
				slog.Debug("Received nil message from claim", "topic", claim.Topic(), "partition", claim.Partition())
				continue
			}

			slog.Debug("Message claimed", "value_len", len(message.Value), "topic", message.Topic, "partition", message.Partition, "offset", message.Offset)
			ctx, err := getTraceCtx(ctx, message.Headers)
			if err != nil {
				slog.Error("Failed to get trace context", "err", err, "topic", message.Topic, "partition", message.Partition)
				session.MarkMessage(message, "") // Mark as processed to avoid blocking partition
			}
			ctx, span := h.tracer.Start(ctx, "kafka-consumer")

			handler, exists := h.handlers[message.Topic]
			if !exists {
				span.SetStatus(codes.Error, "No handler for topic")
				err = fmt.Errorf("no handler for topic %s", message.Topic)
				span.RecordError(err)
				slog.Error("No handler for topic", "topic", message.Topic)
				session.MarkMessage(message, "") // Mark as processed to avoid blocking partition
				continue
			}

			// Create a new context for each handler invocation with a timeout
			// Use session.Context() as the parent, it's cancelled when the session ends.
			handlerCtx, cancel := context.WithTimeout(ctx, h.handlerTimeout)

			err = handler(handlerCtx, message)
			if err != nil {
				span.SetStatus(codes.Error, "Handler error")
				span.RecordError(err, trace.WithAttributes(
					attribute.String("topic", message.Topic),
					attribute.Int64("offset", message.Offset),
				))
				slog.Error("Failed to process sarama msg", "err", err, "topic", message.Topic, "offset", message.Offset)
				// Error handling: depending on the error, you might not want to mark the message.
				// For persistent errors, consider a dead-letter queue strategy.
			}
			cancel() // Ensure cancel is called for the handlerCtx

			session.MarkMessage(message, "") // Mark message as processed
			span.End()
		case <-session.Context().Done(): // Check if the session context is cancelled (e.g., rebalance)
			slog.Info("Session context done, exiting ConsumeClaim", "topic", claim.Topic(), "partition", claim.Partition(), "err", session.Context().Err())
			return session.Context().Err()
		}
	}
}

// Run method updated for Sarama
func (c *Consumer) Run(ctx context.Context) error {
	topics := lo.Keys(c.handlers)
	if len(topics) == 0 {
		return fmt.Errorf("no topics to subscribe to")
	}
	slog.Info("Sarama consumer starting", "topics", topics, "groupID", c.groupID)

	groupHandler := &saramaGroupHandler{
		handlers:       c.handlers,
		handlerTimeout: handlerTimeout,
		tracer:         c.tracer,
	}

	// This loop keeps the consumer active.
	// The `Consume` call will block until the context is cancelled or a non-recoverable error occurs.
	// Sarama handles rebalancing internally when `Consume` is active.
	go func() {
		for err := range c.consumerGroup.Errors() {
			slog.Error("Error from sarama consumer group", "err", err)
		}
	}()

	for {
		// `Consume` should be called in a loop to handle re-joining the group
		// if the consumer leaves for any reason (e.g. network issues, rebalance).
		err := c.consumerGroup.Consume(ctx, topics, groupHandler)
		if err != nil {
			if err == sarama.ErrClosedConsumerGroup {
				slog.Info("Consumer group closed gracefully.")
				return nil
			}
			if err == context.Canceled || err == context.DeadlineExceeded {
				slog.Info("Context cancelled, shutting down consumer.", "err", err)
				return err
			}
			slog.Error("Error from consumer group consume", "err", err, "groupID", c.groupID, "topics", topics)

			select {
			case <-ctx.Done():
				slog.Info("Context done during error handling, exiting consumer run loop.", "err", ctx.Err())
				return ctx.Err()
			case <-time.After(5 * time.Second):
				slog.Info("Retrying to consume after error...")
			}
		}
		if ctx.Err() != nil {
			slog.Info("Context done, exiting consumer run loop.", "err", ctx.Err())
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	if c.consumerGroup != nil {
		slog.Info("Closing sarama consumer group", "groupID", c.groupID)
		return c.consumerGroup.Close()
	}
	return nil
}
