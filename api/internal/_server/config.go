package server

import (
	"fmt"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/larek-tech/diploma/api/internal/api/chat/model"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
)

// Config is the server configuration.
type Config struct {
	HttpPort         int      `yaml:"http_port"`
	BodyLimit        int      `yaml:"body_limit"`
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	wsConfig         websocket.Config
}

// GetAllowedOrigins assembles all allowed origins from slice into string.
func (c *Config) GetAllowedOrigins() string {
	if c.AllowOrigins == nil {
		c.AllowOrigins = []string{"*"}
	}
	return strings.Join(c.AllowOrigins, ",")
}

// GetAllowedMethods assembles all allowed methods from slice into string.
func (c *Config) GetAllowedMethods() string {
	if c.AllowMethods == nil {
		c.AllowMethods = []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
	}
	return strings.Join(c.AllowMethods, ",")
}

// GetAllowedHeaders assembles all allowed headers from slice into string.
func (c *Config) GetAllowedHeaders() string {
	return strings.Join(c.AllowHeaders, ",")
}

// WsConfig returns configuration for server websocket handlers.
func (c *Config) WsConfig() websocket.Config {
	return websocket.Config{
		Origins: c.AllowOrigins,
		RecoverHandler: func(conn *websocket.Conn) {
			if e := recover(); e != nil {
				err := errs.WrapErr(fmt.Errorf("%v", e), "internal error")
				log.Err(err).Msg("ws panic")

				writeErr := conn.WriteJSON(model.SocketMessage{
					Type:   model.TypeError,
					IsLast: true,
					Err:    err,
				})
				if writeErr != nil {
					log.Warn().Err(errs.WrapErr(writeErr)).Msg("failed send recover message")
				}
			}
		},
	}
}
