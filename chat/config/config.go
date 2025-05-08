package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	server "github.com/larek-tech/diploma/chat/internal/_server"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/grpc_client"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config is the application configuration.
type Config struct {
	LogLevel  string            `yaml:"log_level"`
	Server    server.Config     `yaml:"server"`
	Postgres  postgres.Config   `yaml:"postgres"`
	Jaeger    tracing.Config    `yaml:"jaeger"`
	MLService grpcclient.Config `yaml:"ml_service"`
}

// New creates new Config.
func New(path string) (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, errs.WrapErr(err, "create config")
	}
	return cfg, nil
}
