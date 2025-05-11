package server

// Config is the server configuration.
type Config struct {
	GrpcPort   int    `yaml:"grpc_port"`
	Encryption string `yaml:"encryption"`
}
