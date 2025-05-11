package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/larek-tech/diploma/chat/internal/chat/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ChatDao is a model for chat on data layer.
type ChatDao struct {
	Content   []ChatContent `db:"content"`
	ID        uuid.UUID     `db:"id"`
	UserID    int64         `db:"user_id"`
	Title     string        `db:"title"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
}

// ToProto converts data model into protobuf format.
func (c *ChatDao) ToProto() *pb.Chat {
	content := make([]*pb.Content, len(c.Content))
	for idx := range c.Content {
		content[idx] = c.Content[idx].ToProto()
	}

	return &pb.Chat{
		Id:        c.ID.String(),
		UserId:    c.UserID,
		Title:     c.Title,
		Content:   content,
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}

type ChatContent struct {
	Query    QueryDao    `db:"query"`
	Response ResponseDao `db:"response"`
}

// ToProto converts data model into protobuf format.
func (c *ChatContent) ToProto() *pb.Content {
	return &pb.Content{
		Query:    c.Query.ToProto(),
		Response: c.Response.ToProto(),
	}
}

// QueryDao is a model for query on data layer.
type QueryDao struct {
	ID         int64     `db:"id"`
	UserID     int64     `db:"user_id"`
	ChatID     uuid.UUID `db:"chat_id"`
	Content    string    `db:"content"`
	DomainID   int64     `db:"domain_id"`
	SourceIDs  []string  `db:"source_ids"`
	ScenarioID int64     `db:"scenario_id"`
	Metadata   []byte    `db:"metadata"`
	CreatedAt  time.Time `db:"created_at"`
}

// ToProto converts data model into protobuf format.
func (q *QueryDao) ToProto() *pb.Query {
	var (
		domainID   *int64 = nil
		scenarioID *int64 = nil
	)

	if q.DomainID >= 1 {
		domainID = &q.DomainID
	}
	if q.ScenarioID >= 1 {
		scenarioID = &q.ScenarioID
	}

	return &pb.Query{
		Id:         q.ID,
		UserId:     q.UserID,
		ChatId:     q.ChatID.String(),
		Content:    q.Content,
		DomainId:   domainID,
		SourceIds:  q.SourceIDs,
		ScenarioId: scenarioID,
		Metadata:   q.Metadata,
		CreatedAt:  timestamppb.New(q.CreatedAt),
	}
}

// ResponseDao is a model for response on data layer.
type ResponseDao struct {
	ID        int64          `db:"id"`
	QueryID   int64          `db:"query_id"`
	ChatID    uuid.UUID      `db:"chat_id"`
	Content   string         `db:"content"`
	Status    ResponseStatus `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

// ToProto converts data model into protobuf format.
func (r *ResponseDao) ToProto() *pb.Response {
	return &pb.Response{
		Id:        r.ID,
		QueryId:   r.QueryID,
		ChatId:    r.ChatID.String(),
		Content:   r.Content,
		Status:    pb.ResponseStatus(r.Status),
		CreatedAt: timestamppb.New(r.CreatedAt),
		UpdatedAt: timestamppb.New(r.UpdatedAt),
	}
}
