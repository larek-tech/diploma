package server

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/yogenyslav/pkg/interceptor"
	"github.com/yogenyslav/pkg/storage"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Server
type Server struct {
	cfg    Config
	srv    *grpc.Server
	pg     storage.SQLDatabase
	tracer trace.Tracer
}

func New(cfg Config, pg storage.SQLDatabase, tracer trace.Tracer) *Server {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall), logging.WithLogOnEvents(logging.FinishCall),
	}
	srvOpts := []grpc.ServerOption{grpc.ChainUnaryInterceptor(
		logging.UnaryServerInterceptor(interceptor.LoggerInterceptor(), logOpts...),
	)}
	srv := grpc.NewServer(srvOpts...)

	return &Server{
		cfg:    cfg,
		srv:    srv,
		pg:     pg,
		tracer: tracer,
	}
}

func (s *Server) Start() {
	defer s.srv.GracefulStop()
}
