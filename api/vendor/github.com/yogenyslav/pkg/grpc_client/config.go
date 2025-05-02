package grpcclient

// Config holds grpc client configuration.
type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
