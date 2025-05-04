package model

import (
	"time"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DomainDao is a model for domain on data layer.
type DomainDao struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	UserID    int64     `db:"user_id"`
	SourceIDs []int64   `db:"source_ids"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ToProto converts dao model into protobuf format.
func (d *DomainDao) ToProto() *pb.Domain {
	return &pb.Domain{
		Id:        d.ID,
		Title:     d.Title,
		SourceIds: d.SourceIDs,
		CreatedAt: timestamppb.New(d.CreatedAt),
		UpdatedAt: timestamppb.New(d.UpdatedAt),
	}
}
