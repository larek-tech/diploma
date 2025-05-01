// Package grpcclient is a wrapper for grpc client methods.
package grpcclient

import (
	"errors"
	"net"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// ErrCreateClient is an error when failed to establish a connection from grpc client to server.
	ErrCreateClient = errors.New("can't establish a grpc connection")
	// ErrCloseConn is an errors when a grpc connection wasn't be properly closed.
	ErrCloseConn = errors.New("closing grpc connection failed")
)

// GrpcClient holds available methods of grpc client.
type GrpcClient interface {
	Conn() *grpc.ClientConn
	Close() error
}

// GrpcClient is an internal wrap for grpc.ClientConn.
type grpcClient struct {
	conn *grpc.ClientConn
}

// NewGrpcClient creates new GrpcClient.
func NewGrpcClient(cfg *Config, opts ...grpc.DialOption) (GrpcClient, error) {
	var grpcOpts []grpc.DialOption
	grpcOpts = append(grpcOpts, opts...)
	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	conn, err := grpc.NewClient(addr, grpcOpts...)
	if err != nil {
		return nil, errors.Join(ErrCreateClient, err)
	}

	return &grpcClient{conn: conn}, nil
}

// NewGrpcClientWithInsecure creates new GrpcClient with insecure credentials option enabled.
func NewGrpcClientWithInsecure(cfg *Config) (GrpcClient, error) {
	return NewGrpcClient(cfg, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// Close grpc.ClientConn.
func (c *grpcClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return errors.Join(ErrCloseConn, err)
	}
	return nil
}

// Conn GetConn returns grpc.ClientConn.
func (c *grpcClient) Conn() *grpc.ClientConn {
	return c.conn
}
