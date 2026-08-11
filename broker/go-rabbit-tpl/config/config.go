package config

type TRabbitMQ struct {
	URL        string `yaml:"url"`
	Queue      string `yaml:"queue"`
	DLXQueue   string `yaml:"dlx_queue"`
	Exchange   string `yaml:"exchange"`
	RetryCount int    `yaml:"retry_count"`
	RetryDelay int    `yaml:"retry_delay"`
}

type Config struct {
	RabbitMQ TRabbitMQ `yaml:"rabbitmq"`
}
