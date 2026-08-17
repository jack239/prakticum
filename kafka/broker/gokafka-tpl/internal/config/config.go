package config

import "time"

type TKafka struct {
	Brokers    []string      `yaml:"brokers"`
	Topics     []string      `yaml:"topic"`
	GroupID    string        `yaml:"group_id"`
	RetryCount int           `yaml:"retry_count"`
	RetryDelay time.Duration `yaml:"retry_delay"`
}

type Config struct {
	Kafka TKafka `yaml:"kafka"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	return cfg, nil
}
