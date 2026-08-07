package rabbitmq

import (
	"log"

	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	monitor *Monitor
}

func (c *Consumer) ConsumeJSON(queue string, handler func(*models.Message) bool) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(queue, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	msgs, err := ch.Consume(queue, "order-processor", false, false, false, false, nil)
	if err != nil { // проверяем ошибку
		log.Fatal("Ошибка подписки на сообщения:", err) // выходим при ошибке
	}

	log.Println("Ожидаем сообщения о заказах...") // логируем готовность

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

func NewConsumer(conn *amqp.Connection, monitor *Monitor) *Consumer {
	return &Consumer{conn: conn, monitor: monitor}
}
