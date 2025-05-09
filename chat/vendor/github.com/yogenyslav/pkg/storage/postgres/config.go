package postgres

import (
	"fmt"
)

// Config is the configuration for the PostgreSQL client.
type Config struct {
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	DB           string `yaml:"db"`
	Ssl          bool   `yaml:"ssl"`
	Port         int    `yaml:"port"`
	RetryTimeout int    `yaml:"retry_timeout"`
}

// URL assembles config values into a conn string.
func (cfg *Config) URL() string {
	sslMode := "disable"
	if cfg.Ssl {
		sslMode = "enable"
	}
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		sslMode,
	)
}
