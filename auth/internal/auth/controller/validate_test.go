package controller

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	jwtware "github.com/golang-jwt/jwt/v5"
	"github.com/larek-tech/diploma/auth/internal/auth/controller/mocks"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/larek-tech/diploma/auth/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	secret := "test_secret"
	cfg := jwt.Config{
		Secret:     secret,
		Expire:     1,
		Encryption: "1F7006E3A96D34CAB69A15F365FED784",
	}
	provider := jwt.New(cfg)

	validMeta := &pb.UserAuthMetadata{
		UserId: 1,
		Roles:  []int64{1, 2},
	}
	validToken, err := provider.CreateAccessToken(validMeta)
	require.NoError(t, err)

	expiredToken := func() string {
		claims := jwtware.MapClaims{
			"sub":   int64(1),
			"roles": []int64{1, 2},
			"exp":   time.Now().Add(-time.Hour).Unix(),
		}
		token := jwtware.NewWithClaims(jwtware.SigningMethodHS256, claims)
		s, err := token.SignedString([]byte(secret))
		require.NoError(t, err)
		return s
	}()

	malformedToken := func() string {
		claims := jwtware.MapClaims{
			"sub":   "not_a_number",
			"roles": "not_an_array",
		}
		token := jwtware.NewWithClaims(jwtware.SigningMethodHS256, claims)
		s, err := token.SignedString([]byte(secret))
		require.NoError(t, err)
		return s
	}()

	tests := []struct {
		name           string
		token          string
		expectedError  error
		expectedResult *pb.ValidateResponse
	}{
		{
			name:          "ValidatesSuccessfullyWithValidToken",
			token:         validToken,
			expectedError: nil,
			expectedResult: &pb.ValidateResponse{
				Meta: validMeta,
			},
		},
		{
			name:           "FailsWithInvalidToken",
			token:          "invalid_token",
			expectedError:  errors.New("invalid token"),
			expectedResult: nil,
		},
		{
			name:           "FailsWithExpiredToken",
			token:          expiredToken,
			expectedError:  errors.New("expired token"),
			expectedResult: nil,
		},
		{
			name:           "FailsWithMalformedClaims",
			token:          malformedToken,
			expectedError:  base64.CorruptInputError(1),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := new(mocks.MockAuthRepo)
			tracer := noop.NewTracerProvider().Tracer("")
			ctrl := New(tracer, mockRepo, provider)

			req := &pb.ValidateRequest{Token: tt.token}
			resp, err := ctrl.Validate(context.Background(), req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tt.expectedError)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResult.Meta.UserId, resp.Meta.UserId)
				assert.Equal(t, tt.expectedResult.Meta.Roles, resp.Meta.Roles)
			}
		})
	}
}
