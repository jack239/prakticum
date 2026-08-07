package rabbtmq

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn    *amqp.Connection
	monitor *Monitor
}

func NewProducer(conn *amqp.Connection, monitor *Monitor) *Producer {
	return &Producer{conn: conn, monitor: monitor}
}

func (p *Producer) PublishJSON(exchange, routingKey string, msg *Message) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body := json.Marshal(msg)
	args := amqp.Table{"x-max-priority": int32(msg.Priority)}

	q, err := ch.QueueDeclare("alerts.priority", true, false, false, false, args)
	if err != nil {
		p.monitor.error++
		return err
	}
	ch.QueueBind(q.Name, "alert.priority", "alerts.events", false, nil)
	if err != nil {
		p.monitor.error++
		return err
	}

	err = ch.Publish("alerts.events", "alert.priority", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
		Priority:    p, // задаём приоритет сообщения: выше -> раньше доставляется из очереди
	})

	if err != nil {
		p.monitor.error++
		return err
	}

	p.monitor.sent++
	return nil
}
