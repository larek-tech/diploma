package qaas

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.dataddo.com/pgq"
	"go.dataddo.com/pgq/x/schema"
)

type Publisher struct {
	pub pgq.Publisher
	db  *sql.DB
}

func NewPublisher(db *sql.DB) *Publisher {
	return &Publisher{
		pub: pgq.NewPublisher(db),
		db:  db,
	}
}

type PublishOptions struct {
	Queue        Queue
	ScheduledFor *time.Time
	SourceQueue  Queue
	MsgType      string
}

type PublishOption func(*PublishOptions)

func WithMsgType(msgType string) PublishOption {
	return func(opts *PublishOptions) {
		opts.MsgType = msgType
	}
}

func WithScheduledFor(t *time.Time) PublishOption {
	return func(opts *PublishOptions) {
		opts.ScheduledFor = t
	}
}

func WithQueue(queue Queue) PublishOption {
	return func(opts *PublishOptions) {
		opts.Queue = queue
	}
}

func WithSourceQueue(queue Queue) PublishOption {
	return func(opts *PublishOptions) {
		opts.SourceQueue = queue
	}
}

func (p *Publisher) Publish(ctx context.Context, rawMsg []any, opts ...PublishOption) ([]string, error) {
	options := &PublishOptions{}
	for _, opt := range opts {
		opt(options)
	}
	if options.ScheduledFor == nil {
		now := time.Now()
		options.ScheduledFor = &now
	}
	if options.Queue == "" {
		return nil, fmt.Errorf("queue name is required")
	}

	msgs := make([]*pgq.MessageOutgoing, len(rawMsg))
	for i := 0; i < len(rawMsg); i++ {
		switch v := rawMsg[i].(type) {
		case SiteJob, PageJob, EmbedJob:
			payload, err := json.Marshal(rawMsg[i])
			if err != nil {
				return nil, fmt.Errorf("failed to marshal message: %w", err)
			}
			// TODO: move metadata keys as constants
			msgs[i] = &pgq.MessageOutgoing{
				ScheduledFor: options.ScheduledFor,
				Payload:      payload,
				Metadata: pgq.Metadata{
					"type":        options.MsgType,
					"queue":       string(options.Queue),
					"objType":     reflect.TypeOf(v).Name(),
					"sourceQueue": string(options.SourceQueue),
				},
			}
		default:
			return nil, fmt.Errorf("unsupported message type: %T", rawMsg[i])
		}

	}

	msgIDs, err := p.pub.Publish(ctx, string(options.Queue), msgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}
	return lo.Map(msgIDs, func(v uuid.UUID, _ int) string { return v.String() }), nil
}

func (p Publisher) CreateAllTables(queues []Queue) error {
	for _, queue := range queues {
		query := schema.GenerateCreateTableQuery(string(queue))
		if _, err := p.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table for queue %s: %w", queue, err)
		}
	}
	return nil
}
