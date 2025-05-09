// Package postgres provides a simple wrapper around the pgx library to interact with a PostgreSQL database.
package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type contextKey uint8

const (
	// TxKey is a key for the transaction stored in context.
	TxKey contextKey = iota
)

var (
	// ErrCreatePostgres is an error when a connection to postgres can't be established.
	ErrCreatePostgres = errors.New("failed to connect to postgres")
	// ErrNoTxFound is an error when no transaction was found, but tried to access it.
	ErrNoTxFound = errors.New("no transaction found")
)

// Postgres wraps pgxpool.Pool and adds tracer to all operations.
type Postgres struct {
	pool   *pgxpool.Pool
	tracer trace.Tracer
}

// New creates a new Postgres instance.
func New(cfg *Config, tracer trace.Tracer) (Postgres, error) {
	pgConfig, err := pgxpool.ParseConfig(cfg.URL())
	if err != nil {
		return Postgres{}, fmt.Errorf("parse postgres connection string: %w", err)
	}

	// set few connection options
	pgConfig.ConnConfig.ConnectTimeout = time.Second * time.Duration(cfg.RetryTimeout)
	pgConfig.ConnConfig.ValidateConnect = func(ctx context.Context, conn *pgconn.PgConn) error {
		return conn.Ping(ctx)
	}

	// create pgx pool for postgres
	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return Postgres{}, fmt.Errorf("connect to postgres: %w", err)
	}

	return Postgres{
		pool:   pool,
		tracer: tracer,
	}, nil
}

// GetPool returns the underlying pgxpool.Pool.
func (pg Postgres) GetPool() *pgxpool.Pool {
	return pg.pool
}

// Close closes the underlying db connection.
func (pg Postgres) Close() {
	pg.pool.Close()
}

// GetTx returns a transaction from ctx or an error if there is no tx.
func (pg Postgres) GetTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("get transaction: %w", ErrNoTxFound)
	}
	return tx, nil
}

// BeginSerializable starts a new transaction with serializable isolation level.
func (pg Postgres) BeginSerializable(ctx context.Context) (context.Context, error) {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(ctx, "Postgres.BeginSerializable")
		defer span.End()
	}

	tx, err := pg.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return ctx, fmt.Errorf("starting a serializable tx failed: %w", err)
	}

	return context.WithValue(ctx, TxKey, tx), nil
}

// CommitTx commits the transaction.
func (pg Postgres) CommitTx(ctx context.Context) error {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(ctx, "Postgres.CommitTx")
		defer span.End()
	}

	tx, err := pg.GetTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// RollbackTx rolls back the transaction.
func (pg Postgres) RollbackTx(ctx context.Context) error {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(ctx, "Postgres.RollbackTx")
		defer span.End()
	}

	tx, err := pg.GetTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := tx.Rollback(ctx); err != nil {
		if !errors.Is(err, pgx.ErrTxClosed) {
			return fmt.Errorf("failed to rollback transaction: %w", err)
		}
	}
	return nil
}

// Query executes a query that returns a single row.
func (pg Postgres) Query(ctx context.Context, dest any, query string, args ...any) error {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(
			ctx,
			"Postgres.Query",
			trace.WithAttributes(attribute.String("query", query)),
		)
		defer span.End()
	}

	if err := pgxscan.Get(ctx, pg.pool, dest, query, args...); err != nil {
		return fmt.Errorf("failed to get row: %w", err)
	}
	return nil
}

// QuerySlice executes a query that returns multiple rows.
func (pg Postgres) QuerySlice(ctx context.Context, dest any, query string, args ...any) error {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(
			ctx,
			"Postgres.QuerySlice",
			trace.WithAttributes(attribute.String("query", query)),
		)
		defer span.End()
	}

	if err := pgxscan.Select(ctx, pg.pool, dest, query, args...); err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	return nil
}

// Exec executes a query that doesn't return any rows.
func (pg Postgres) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(
			ctx,
			"Postgres.Exec",
			trace.WithAttributes(attribute.String("query", query)),
		)
		defer span.End()
	}

	tag, err := pg.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to exec: %w", err)
	}
	return tag.RowsAffected(), nil
}

// QueryTx executes a query that returns a single row in a transaction.
func (pg Postgres) QueryTx(ctx context.Context, dest any, query string, args ...any) error {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(
			ctx,
			"Postgres.QueryTx",
			trace.WithAttributes(attribute.String("query", query)),
		)
		defer span.End()
	}

	tx, err := pg.GetTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := pgxscan.Get(ctx, tx, dest, query, args...); err != nil {
		return fmt.Errorf("failed to get row in transaction: %w", err)
	}
	return nil
}

// QuerySliceTx executes a query that returns multiple rows in a transaction.
func (pg Postgres) QuerySliceTx(ctx context.Context, dest any, query string, args ...any) error { //nolint:dupl // because of pgx api
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(
			ctx,
			"Postgres.QuerySliceTx",
			trace.WithAttributes(attribute.String("query", query)),
		)
		defer span.End()
	}

	tx, err := pg.GetTx(ctx)
	if err != nil {
		return fmt.Errorf("get transaction: %w", err)
	}

	if err := pgxscan.Select(ctx, tx, dest, query, args...); err != nil {
		return fmt.Errorf("failed to get rows in transaction: %w", err)
	}
	return nil
}

// ExecTx executes a query that doesn't return any rows in a transaction.
func (pg Postgres) ExecTx(ctx context.Context, query string, args ...any) (int64, error) {
	if pg.tracer != nil {
		var span trace.Span
		ctx, span = pg.tracer.Start(
			ctx,
			"Postgres.ExecTx",
			trace.WithAttributes(attribute.String("query", query)),
		)
		defer span.End()
	}

	tx, err := pg.GetTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("get transaction: %w", err)
	}

	tag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to exec in transaction: %w", err)
	}
	return tag.RowsAffected(), nil
}
