package model

import (
	"time"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// RoleDao is a role model on data layer.
type RoleDao struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// ToProto converts data model into protobuf format.
func (r *RoleDao) ToProto() *pb.Role {
	return &pb.Role{
		Id:        r.ID,
		Name:      r.Name,
		CreatedAt: timestamppb.New(r.CreatedAt),
	}
}
