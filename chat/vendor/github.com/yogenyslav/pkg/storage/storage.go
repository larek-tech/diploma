// Package storage provides wrappers around different storage solutions.
package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SQLDatabase is an interface that wraps the basic SQL operations.
type SQLDatabase interface {
	// BeginSerializable starts a new transaction with serializable isolation level.
	BeginSerializable(ctx context.Context) (context.Context, error)
	// GetTx returns a transaction from ctx or an error if there is no tx.
	GetTx(ctx context.Context) (pgx.Tx, error)
	// CommitTx commits the transaction.
	CommitTx(ctx context.Context) error
	// RollbackTx rolls back the transaction.
	RollbackTx(ctx context.Context) error
	// Query executes a query that returns a single row.
	Query(ctx context.Context, dest any, query string, args ...any) error
	// QuerySlice executes a query that returns multiple rows.
	QuerySlice(ctx context.Context, dest any, query string, args ...any) error
	// Exec executes a query that doesn't return any rows.
	// Returns number of affected rows.
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	// QueryTx executes a query that returns a single row in a transaction.
	QueryTx(ctx context.Context, dest any, query string, args ...any) error
	// QuerySliceTx executes a query that returns multiple rows in a transaction.
	QuerySliceTx(ctx context.Context, dest any, query string, args ...any) error
	// ExecTx executes a query that doesn't return any rows in a transaction.
	// Returns number of affected rows.
	ExecTx(ctx context.Context, query string, args ...any) (int64, error)
	// Close closes the database connection.
	Close()
}

// MongoDatabase is an interface that wraps the basic MongoDB operations.
type MongoDatabase interface {
	// InsertOne inserts a single document into the collection.
	InsertOne(ctx context.Context, coll string, doc interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	// FindOne finds a single document in the collection.
	FindOne(ctx context.Context, coll string, filter, dest interface{}, opts ...*options.FindOneOptions) error
	// FindMany finds multiple documents in the collection.
	FindMany(ctx context.Context, coll string, filter, dest interface{}, opts ...*options.FindOptions) error
	// UpdateOne updates a single document in the collection.
	UpdateOne(ctx context.Context, coll string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// UpdateMany updates multiple documents in the collection.
	UpdateMany(ctx context.Context, coll string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	// DeleteOne deletes a single document from the collection.
	DeleteOne(ctx context.Context, coll string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	// Close closes the MongoDB client.
	Close() error
}

// Cache is an interface that wraps the basic cache operations.
type Cache interface {
	// SetStruct sets a struct in the cache.
	SetStruct(ctx context.Context, k string, v any, exp time.Duration) error
	// SetPrimitive sets a primitive in the cache.
	SetPrimitive(ctx context.Context, k string, v any, exp time.Duration) error
	// GetStruct gets a struct from the cache.
	GetStruct(ctx context.Context, dest any, k string) error
	// GetString gets a string from the cache.
	GetString(ctx context.Context, k string) (string, error)
	// GetInt gets an int from the cache.
	GetInt(ctx context.Context, k string) (int, error)
	// GetInt64 gets an int64 from the cache.
	GetInt64(ctx context.Context, k string) (int64, error)
	// GetFloat gets a float64 from the cache.
	GetFloat(ctx context.Context, k string) (float64, error)
	// GetBool gets a bool from the cache.
	GetBool(ctx context.Context, k string) (bool, error)
	// GetBytes gets a byte slice from the cache.
	GetBytes(ctx context.Context, k string) ([]byte, error)
	// Del deletes a key from the cache.
	Del(ctx context.Context, k string) error
}
