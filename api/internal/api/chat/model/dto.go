package model

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
)

type SocketMessageType string

const (
	TypeAuth  SocketMessageType = "auth"
	TypeQuery SocketMessageType = "query"
	TypeChunk SocketMessageType = "chunk"
	TypeError SocketMessageType = "error"
)

// SocketMessage is a model for incoming and outgoing messages for websocket.
type SocketMessage struct {
	Type          SocketMessageType `json:"type"`
	Content       string            `json:"content"`
	IsChunked     bool              `json:"isChunked"`
	IsLast        bool              `json:"isLast"`
	SourceIDs     []string          `json:"sourceIDs"`
	QueryMetadata QueryMetadata     `json:"queryMetadata,omitempty"`
	Err           error             `json:"error,omitempty"`
}

// QueryMetadata stores information about chosen domain, sources and scenario for query.
type QueryMetadata struct {
	DomainID *int64       `json:"domainID"`
	Scenario *pb.Scenario `json:"scenario"`
}
