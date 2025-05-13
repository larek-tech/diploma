package model

// SocketMessageType enum defining type of socket message.
type SocketMessageType string

const (
	// TypeAuth content is auth token (without Bearer).
	TypeAuth SocketMessageType = "auth"
	// TypeQuery content is general message.
	TypeQuery SocketMessageType = "query"
	// TypeChunk content is chunked response.
	TypeChunk SocketMessageType = "chunk"
	// TypeError content is empty, got error, chat is stopped.
	TypeError SocketMessageType = "error"
)

// SocketMessage is a model for incoming and outgoing messages for websocket.
type SocketMessage struct {
	Type       SocketMessageType `json:"type"`
	Content    string            `json:"content,omitempty"`
	IsChunked  bool              `json:"isChunked"`
	IsLast     bool              `json:"isLast"`
	DomainID   int64             `json:"domainID"`
	ScenarioID int64             `json:"scenarioID"`
	Err        string            `json:"error,omitempty"`
}
