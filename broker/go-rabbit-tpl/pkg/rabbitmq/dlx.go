package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func SetupDLX(ch *amqp.Channel, queueName, dlxQueueName, dlxExchange string) error {
	if err := ch.ExchangeDeclare(dlxExchange, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("payments.main", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(queueName, dlxQueueName, dlxExchange, false, nil); err != nil {
		return err
	}

	args := amqp.Table{ // аргументы очереди задержки
		"x-message-ttl":             int32(30_000), // хранить сообщение 30 секунд
		"x-dead-letter-exchange":    dlxExchange,   // по истечении TTL отправить в DLX
		"x-dead-letter-routing-key": dlxQueueName,  // ключ для возврата в основную
	}
	if _, err := ch.QueueDeclare("payments.retry.30s", true, false, false, false, args); err != nil { // объявляем retry-очередь
		return err
	}
	return nil
}
