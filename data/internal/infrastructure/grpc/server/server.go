package server

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

const port = 50051

type Server struct {
	srv *grpc.Server
}

func New() *Server {
	srvOpts := []grpc.ServerOption{}
	srv := grpc.NewServer(srvOpts...)
	return &Server{srv: srv}
}

func (s Server) Run() error {
	defer s.srv.GracefulStop()

	slog.Info("Starting server")

	lis, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("create net listener: %w", err)
	}
	slog.Info("Listening grpc on port", "port", port)

	if err = s.srv.Serve(lis); err != nil {
		return fmt.Errorf("shutting down: %w", err)
	}
	return nil
}

func (s Server) GetSrv() *grpc.Server {
	return s.srv
}
