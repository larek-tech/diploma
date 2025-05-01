package server

// Config is the server configuration.
type Config struct {
	GrpcPort int    `yaml:"grpc_port"`
	LogLevel string `yaml:"log_level"`
}
