package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/larek-tech/diploma/auth/internal/auth/controller/mocks"
	"github.com/larek-tech/diploma/auth/internal/auth/model"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/larek-tech/diploma/auth/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yogenyslav/pkg/secure"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	errDBError := errors.New("db error")

	tests := []struct {
		name           string
		setupMocks     func(mockRepo *mocks.MockAuthRepo)
		request        *pb.LoginRequest
		expectedError  error
		expectedResult *pb.LoginResponse
	}{
		{
			name: "LoginSuccessWithValidCredentials",
			setupMocks: func(mockRepo *mocks.MockAuthRepo) {
				password := "test123456"
				hashedPassword, _ := secure.HashPassword(password)
				user := model.UserDao{
					ID:           1,
					Email:        "test@test.com",
					HashPassword: hashedPassword,
				}
				roles := []int64{1, 2}

				mockRepo.On("FindOneByEmail", mock.Anything, user.Email).Return(user, nil)
				mockRepo.On("FindUserRoles", mock.Anything, user.ID).Return(roles, nil)
			},
			request: &pb.LoginRequest{
				Email:    "test@test.com",
				Password: "test123456",
			},
			expectedError: nil,
			expectedResult: &pb.LoginResponse{
				Token: "mocked_token", // Replace with actual token if needed
				Type:  jwt.TypeBearerToken,
				Meta: &pb.UserAuthMetadata{
					UserId: 1,
					Roles:  []int64{1, 2},
				},
			},
		},
		{
			name: "LoginFailsWithInvalidPassword",
			setupMocks: func(mockRepo *mocks.MockAuthRepo) {
				hashedPassword, _ := secure.HashPassword("correctpass")
				user := model.UserDao{
					ID:           1,
					Email:        "test@test.com",
					HashPassword: hashedPassword,
				}
				mockRepo.On("FindOneByEmail", mock.Anything, user.Email).Return(user, nil)
			},
			request: &pb.LoginRequest{
				Email:    "test@test.com",
				Password: "wrongpass",
			},
			expectedError:  ErrInvalidCredentials,
			expectedResult: nil,
		},
		{
			name: "LoginFailsWithNonExistentUser",
			setupMocks: func(mockRepo *mocks.MockAuthRepo) {
				mockRepo.On("FindOneByEmail", mock.Anything, "nonexistent@test.com").Return(model.UserDao{}, ErrUserNotFound)
			},
			request: &pb.LoginRequest{
				Email:    "nonexistent@test.com",
				Password: "anypass",
			},
			expectedError:  ErrUserNotFound,
			expectedResult: nil,
		},
		{
			name: "LoginFailsWhenRolesLookupFails",
			setupMocks: func(mockRepo *mocks.MockAuthRepo) {
				password := "test123456"
				hashedPassword, _ := secure.HashPassword(password)
				user := model.UserDao{
					ID:           1,
					Email:        "test@test.com",
					HashPassword: hashedPassword,
				}
				mockRepo.On("FindOneByEmail", mock.Anything, user.Email).Return(user, nil)
				mockRepo.On("FindUserRoles", mock.Anything, user.ID).Return([]int64{}, errDBError)
			},
			request: &pb.LoginRequest{
				Email:    "test@test.com",
				Password: "test123456",
			},
			expectedError:  errDBError,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.MockAuthRepo)
			mockJWT := &jwt.Provider{}
			tracer := noop.NewTracerProvider().Tracer("")

			ctrl := New(tracer, mockRepo, mockJWT)

			tt.setupMocks(mockRepo)

			resp, err := ctrl.Login(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResult.Type, resp.Type)
				assert.Equal(t, tt.expectedResult.Meta.UserId, resp.Meta.UserId)
				assert.Equal(t, tt.expectedResult.Meta.Roles, resp.Meta.Roles)
			}
		})
	}
}
