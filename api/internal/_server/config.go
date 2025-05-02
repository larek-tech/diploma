package server

import (
	"strings"
)

// Config is the server configuration.
type Config struct {
	HttpPort         int      `yaml:"http_port"`
	BodyLimit        int      `yaml:"body_limit"`
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
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
