package tracing

import (
	"net"
	"strconv"
)

// Config is the struct for tracing configuration.
type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// URL returns the URL for the tracing server.
func (c *Config) URL() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
