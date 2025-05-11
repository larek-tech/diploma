package model

import (
	"time"

	"github.com/larek-tech/diploma/domain/internal/domain/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserDao is a data model for user on data layer.
type UserDao struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	HashPassword string    `db:"hash_password"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	IsDeleted    bool      `db:"is_deleted"`
}

// ToProto converts data model into protobuf format.
func (u *UserDao) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
