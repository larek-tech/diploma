package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	server "github.com/larek-tech/diploma/auth/internal/_server"
	"github.com/larek-tech/diploma/auth/pkg/jwt"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config is the application configuration.
type Config struct {
	LogLevel string          `yaml:"log_level"`
	Server   server.Config   `yaml:"server"`
	Postgres postgres.Config `yaml:"postgres"`
	Jaeger   tracing.Config  `yaml:"jaeger"`
	Jwt      jwt.Config      `yaml:"jwt"`
}

// New creates new Config.
func New(path string) (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, errs.WrapErr(err, "create config")
	}
	return cfg, nil
}
