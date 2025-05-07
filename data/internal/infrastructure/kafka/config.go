package kafka

type Config struct {
	BootstrapServers string `env:"KAFKA_SERVERS"`
}
