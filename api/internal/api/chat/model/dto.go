package model

import (
	"github.com/larek-tech/diploma/api/internal/domain/pb"
)

// SocketMessageType enum defining type of socket message.
type SocketMessageType string

const (
	// TypeAuth content is auth token (without Bearer).
	TypeAuth  SocketMessageType = "auth"
	// TypeQuery content is general message.
	TypeQuery SocketMessageType = "query"
	// TypeChunk content is chunked response.
	TypeChunk SocketMessageType = "chunk"
	// TypeError content is empty, got error, chat is stopped.
	TypeError SocketMessageType = "error"
)

// SocketMessage is a model for incoming and outgoing messages for websocket.
type SocketMessage struct {
	Type          SocketMessageType `json:"type"`
	Content       string            `json:"content,omitempty"`
	IsChunked     bool              `json:"isChunked"`
	IsLast        bool              `json:"isLast"`
	SourceIDs     []string          `json:"sourceIDs,omitempty"`
	QueryMetadata QueryMetadata     `json:"queryMetadata,omitempty"`
	Err           string            `json:"error,omitempty"`
}

// QueryMetadata stores information about chosen domain, sources and scenario for query.
type QueryMetadata struct {
	DomainID *int64       `json:"domainID,omitempty"`
	Scenario *pb.Scenario `json:"scenario,omitempty"`
}
