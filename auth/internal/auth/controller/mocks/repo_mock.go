package mocks

import (
	"context"

	"github.com/larek-tech/diploma/auth/internal/auth/model"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) FindOneByEmail(ctx context.Context, email string) (model.UserDao, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.UserDao), args.Error(1)
}

func (m *MockAuthRepo) FindUserRoles(ctx context.Context, userID int64) ([]int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]int64), args.Error(1)
}
