package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/interceptor"
	"github.com/yogenyslav/pkg/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Server holds grpc auth server and its dependencies.
type Server struct {
	cfg    Config
	srv    *grpc.Server
	pg     storage.SQLDatabase
	tracer trace.Tracer
}

// New creates new Server.
func New(cfg Config, pg storage.SQLDatabase, tracer trace.Tracer) *Server {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall), logging.WithLogOnEvents(logging.FinishCall),
	}
	srvOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptor.LoggerInterceptor(), logOpts...),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	srv := grpc.NewServer(srvOpts...)

	return &Server{
		cfg:    cfg,
		srv:    srv,
		pg:     pg,
		tracer: tracer,
	}
}

func (s *Server) GetSrv() *grpc.Server {
	return s.srv
}

// Start starts listening gRPC requests on port.
func (s *Server) Start() {
	defer s.srv.GracefulStop()

	log.Info().Int("port", s.cfg.GrpcPort).Msg("Starting gRPC server")

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go s.listenGrpc(errCh)

	select {
	case <-stopCh:
		log.Info().Msg("Server graceful shutdown")
	case err := <-errCh:
		log.Err(err).Msg("Server fatal error")
	}
}

func (s *Server) listenGrpc(errCh chan error) {
	lis, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(s.cfg.GrpcPort)))
	if err != nil {
		errCh <- errs.WrapErr(err, "create net listener")
		return
	}

	if err = s.srv.Serve(lis); err != nil {
		errCh <- errs.WrapErr(err)
		return
	}
}
