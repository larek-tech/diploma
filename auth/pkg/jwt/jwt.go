package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/larek-tech/diploma/auth/internal/auth/pb"
	"github.com/yogenyslav/pkg/secure"
)

const (
	// TypeBearerToken value "Bearer" for the token type field.
	TypeBearerToken string = "Bearer"
)

var (
	// ErrJwtSignMethod is an error when jwt signing method is wrong.
	ErrJwtSignMethod = errors.New("unexpected signing method")
)

// Config is a config for jwt module.
type Config struct {
	Secret     string `yaml:"secret"`
	Expire     int    `yaml:"expire"`
	Encryption string `yaml:"encryption"`
}

// Provider implements jwt token generation and validation.
type Provider struct {
	cfg Config
}

func New(cfg Config) *Provider {
	return &Provider{
		cfg: cfg,
	}
}

func (j *Provider) CreateAccessToken(meta *pb.UserAuthMetadata) (string, error) {
	key := []byte(j.cfg.Secret)

	jwtClaims := jwt.MapClaims{
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(j.cfg.Expire))),
		"sub":   meta.GetUserId(),
		"roles": meta.GetRoles(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := accessToken.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	if j.cfg.Encryption != "" {
		return secure.Encrypt(signedToken, j.cfg.Encryption)
	}

	return signedToken, nil
}

func (j *Provider) ParseAccessToken(accessTokenString string) (*jwt.Token, error) {
	var err error

	if j.cfg.Encryption != "" {
		accessTokenString, err = secure.Decrypt(accessTokenString, j.cfg.Encryption)
		if err != nil {
			return nil, fmt.Errorf("decrypt token: %w", err)
		}
	}

	accessToken, err := jwt.Parse(accessTokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse token: %w", ErrJwtSignMethod)
		}
		return []byte(j.cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	return accessToken, nil
}
