package rabbtmq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	monitor *Monitor
}

func (c *Consumer) ConsumeJSON(queue string, handler func(*Message) bool) error {
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

	// Обрабатываем входящие сообщения
	for d := range msgs { // читаем поток сообщений
		var message Message
		if err := json.Unmarshal(d.Body, &message); err != nil {
			c.monitor.error++
			continue
		}
		if err := handler(&message); err != nil {
			c.monitor.error++
			continue
		}

	}
}
