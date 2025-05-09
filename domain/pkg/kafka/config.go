package kafka

// Config holds the list of brokers and topics.
type Config struct {
	Brokers        []Broker `yaml:"brokers"`
	Topics         []Topic  `yaml:"topics"`
	OffsetNewest   bool     `yaml:"offset_newest"`
	CommitInterval int      `yaml:"commit_interval"`
}

// Broker is the struct for Kafka broker.
type Broker struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Topic is the struct for Kafka topic.
type Topic struct {
	Name       string `yaml:"name"`
	Partitions int32  `yaml:"partitions"`
}
