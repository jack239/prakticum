package client

import "github.com/segmentio/kafka-go"

type KafkaClient struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func (client *KafkaClient) GetWriter() *kafka.Writer {
	return client.writer
}

func (client *KafkaClient) GetReader() *kafka.Reader {
	return client.reader
}

func (client *KafkaClient) Close() error {
	if err := client.writer.Close(); err != nil {
		return err
	}
	if err := client.reader.Close(); err != nil {
		return err
	}
	return nil
}

func NewKafkaClient(
	brokers []string,
	topic string,
	groupID string,
) (*KafkaClient, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{}, // Хеширование по ключу
	}
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: 0, // Отключаем авто-коммит
		StartOffset:    kafka.LastOffset,
	})
	return &KafkaClient{
		writer: writer,
		reader: reader,
	}, nil
}
