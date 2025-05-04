package server

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

func (s Server) Run() {
	defer s.srv.GracefulStop()

	slog.Info("Starting server")

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	go s.listenGrpc(errCh)

	select {
	case <-stopCh:
		slog.Info("Received shutdown signal")
	case err := <-errCh:

		slog.Error("Error occurred", "error", err)
	}

}

func (s Server) GetSrv() *grpc.Server {
	return s.srv
}

func (s Server) listenGrpc(errCh chan error) {
	lis, err := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(port)))
	if err != nil {
		errCh <- fmt.Errorf("create net listener: %w", err)
		return
	}
	slog.Info("Listening grpc on port", "port", port)

	if err = s.srv.Serve(lis); err != nil {
		errCh <- fmt.Errorf("shutting down: %w", err)
		return
	}
}
