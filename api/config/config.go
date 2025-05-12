package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	server "github.com/larek-tech/diploma/api/internal/_server"
	"github.com/yogenyslav/pkg/errs"
	grpcclient "github.com/yogenyslav/pkg/grpc_client"
	"github.com/yogenyslav/pkg/infrastructure/tracing"
	"github.com/yogenyslav/pkg/storage/postgres"
)

// Config is the application configuration.
type Config struct {
	LogLevel      string            `yaml:"log_level"`
	Server        server.Config     `yaml:"server"`
	Postgres      postgres.Config   `yaml:"postgres"`
	Jaeger        tracing.Config    `yaml:"jaeger"`
	AuthService   grpcclient.Config `yaml:"auth_service"`
	DomainService grpcclient.Config `yaml:"domain_service"`
	ChatService   grpcclient.Config `yaml:"chat_service"`
	MLService     grpcclient.Config `yaml:"ml_service"`
}

// New creates new Config.
func New(path string) (Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, errs.WrapErr(err, "create config")
	}
	return cfg, nil
}
