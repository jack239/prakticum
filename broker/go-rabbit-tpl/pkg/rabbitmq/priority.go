package rabbtmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclarePriorityQueue(ch *amqp.Channel, queueName string, maxPriotiy int) error {
	args := amqp.Table{"x-max-priority": int32(maxPriotiy)}
	_, err := ch.QueueDeclare(queueName, true, false, false, false, args)
	return err
}
