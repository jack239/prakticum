package gokafka

import (
	"time"

	"github.com/lekan-pvp/grade/go-kafka/internal/message"
	"github.com/lekan-pvp/grade/go-kafka/internal/monitor"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer     *kafka.Writer
	mon        *monitor.Monitor
	maxRetries int
	retryDelay time.Duration
}

func NewProducer(writer *kafka.Writer, mon *monitor.Monitor, retryCount int, retryDelay time.Duration) *Producer {
	return &Producer{
		writer: writer,
		mon:    mon,
	}
}

func (p *Producer) PublishJSON(topic string, mes *message.Message, partition ...int) error {
	data, err := mes.ToJSON()
	if err != nil {
		p.mon.IncError()
		return err
	}
	kmessage := kafka.Message{
		Key:   []byte(mes.Key),
		Value: data,
	}

	if len(partition) != 0 && partition[0] != 0 {
		kmessage.Partition = partition[0]
	}

	for attempt := range p.maxRetries {
		// err := p.writer.WriteMessages(kmessage)
		if err != nil {
			p.mon.IncError()
			if attempt+1 == p.maxRetries {
				return err
			}
			time.Sleep(p.retryDelay)
			continue
		}
		p.mon.IncPublished()
		break
	}
	return nil
}
