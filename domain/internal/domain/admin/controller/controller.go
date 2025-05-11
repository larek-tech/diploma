package controller

import (
	"context"

	"github.com/larek-tech/diploma/domain/internal/domain/admin/model"
)

type adminRepo interface {
	InsertUser(ctx context.Context, u model.UserDao) (int64, error)
	GetUser(ctx context.Context, id int64) (model.UserDao, error)
	UpdateUser(ctx context.Context, u model.UserDao) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, offset, limit uint64) ([]model.UserDao, error)
}
