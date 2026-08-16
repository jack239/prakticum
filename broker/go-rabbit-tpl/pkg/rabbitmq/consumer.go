package rabbitmq

import (
	"log"

	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
)

type Consumer struct {
	connection *Connection
	monitor    *Monitor
}

func (c *Consumer) ConsumeJSON(queue string, handler func(*models.Message) bool) error {
	msgs, err := c.connection.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range msgs { // читаем поток сообщений
			msg, err := models.FromJSON(delivery.Body)
			if err != nil {
				c.monitor.IncError()
				log.Printf("Failed to parse JSON message: %v", err)
				_ = delivery.Reject(false)
				continue
			}

			if handler(msg) {
				err = delivery.Ack(false)
				if err != nil {
					log.Printf("Failed to send ACK: %v", err)
				} else {
					c.monitor.IncReceived()
				}
			} else {
				err = delivery.Reject(false)
				if err != nil {
					log.Printf("Failed to send NASK: %v", err)
				}
			}

		}
	}()
	return nil
}

func NewConsumer(conn *Connection, monitor *Monitor) *Consumer {
	return &Consumer{connection: conn, monitor: monitor}
}
