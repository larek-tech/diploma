package pkg

import (
	"github.com/larek-tech/diploma/auth/config"
	server "github.com/larek-tech/diploma/auth/internal/_server"
	"github.com/yogenyslav/pkg/errs"
)

const (
	configPath = "./config/config.yaml"
)

func Run() error {
	cfg, err := config.New(configPath)
	if err != nil {
		return errs.WrapErr(err)
	}

	srv := server.New(cfg.Server)
	srv.Start()

	return nil
}
