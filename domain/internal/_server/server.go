package server

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/pkg/errs"
	"github.com/yogenyslav/pkg/interceptor"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Server holds grpc auth server.
type Server struct {
	cfg Config
	srv *grpc.Server
}

// New creates new Server.
func New(cfg Config) *Server {
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
		cfg: cfg,
		srv: srv,
	}
}

// GetSrv returns underlying gRPC conn.
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
