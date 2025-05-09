package server

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	_ "github.com/larek-tech/diploma/api/docs"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/response"
)

// Server holds HTTP api server.
type Server struct {
	cfg Config
	srv *fiber.App
}

// New creates new Server.
func New(cfg Config) *Server {
	srv := fiber.New(fiber.Config{
		BodyLimit:    cfg.BodyLimit,
		ErrorHandler: response.NewErrorHandler(errMap).Handler(&log.Logger),
		AppName:      "Diploma API",
	})

	srv.Use(otelfiber.Middleware())
	srv.Use(logger.New(logger.Config{
		TimeZone: "Europe/Moscow",
	}))
	srv.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.GetAllowedOrigins(),
		AllowMethods:     cfg.GetAllowedMethods(),
		AllowHeaders:     cfg.GetAllowedHeaders(),
		AllowCredentials: cfg.AllowCredentials,
	}))
	srv.Use("/api/v1/swagger/*", swagger.HandlerDefault)
	srv.Use(recovermw.New())

	return &Server{
		cfg: cfg,
		srv: srv,
	}
}

// GetSrv returns underlying HTTP conn.
func (s *Server) GetSrv() *fiber.App {
	return s.srv
}

// Start starts listening HTTP requests on port.
func (s *Server) Start() {
	defer func() {
		if err := s.srv.Shutdown(); err != nil {
			log.Warn().Err(err).Msg("server graceful shutdown failed")
		}
	}()

	log.Info().Int("port", s.cfg.HttpPort).Msg("Starting HTTP server")

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go s.listenHTTP(errCh)

	select {
	case <-stopCh:
		log.Info().Msg("Server graceful shutdown")
	case err := <-errCh:
		log.Err(err).Msg("Server fatal error")
	}
}

func (s *Server) listenHTTP(errCh chan error) {
	addr := net.JoinHostPort("", strconv.Itoa(s.cfg.HttpPort))
	if err := s.srv.Listen(addr); err != nil {
		errCh <- errs.WrapErr(err)
	}
}
