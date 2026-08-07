package rabbitmq

import (
	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn    *Connection
	monitor *Monitor
}

func NewProducer(conn *Connection, monitor *Monitor) *Producer {
	return &Producer{conn: conn, monitor: monitor}
}

func (p *Producer) PublishJSON(exchange, routingKey string, msg *models.Message) error {
	body, err := msg.ToJSON()
	if err != nil {
		p.monitor.IncError()
		return err
	}

	err = p.conn.Channel().Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
		Priority:    uint8(msg.Priority), // задаём приоритет сообщения: выше -> раньше доставляется из очереди
	})
	if err != nil {
		p.monitor.IncError()
		return err
	}

	p.monitor.IncSent()
	return nil
}
