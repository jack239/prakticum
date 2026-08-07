package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Подключение к RabbitMQ серверу
	// Формат URL: amqp://username:password@host:port/vhost
	log.Println("Подключаемся к RabbitMQ...")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("Успешно подключились к RabbitMQ")

	// Создание канала для работы с очередями и сообщениями
	// Канал — это виртуальное соединение внутри TCP-соединения
	log.Println("Создаем канал...")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка создания канала: %v", err)
	}
	defer ch.Close()
	log.Println("Канал создан успешно")

	// Объявление очереди
	// Параметры: name, durable, autoDelete, exclusive, noWait, args
	// durable=true — очередь переживет перезапуск сервера
	// autoDelete=false — очередь не удалится когда отключится последний потребитель
	log.Println("Объявляем очередь 'demo.q'...")
	q, err := ch.QueueDeclare(
		"demo.q", // имя очереди
		true,     // durable — очередь будет сохранена на диск
		false,    // autoDelete — не удалять при отсутствии потребителей
		false,    // exclusive — очередь доступна другим соединениям
		false,    // noWait — ждать подтверждения от сервера
		nil,      // args — дополнительные аргументы
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %v", err)
	}
	log.Printf("Очередь '%s' объявлена успешно", q.Name)

	// Создание контекста с таймаутом для публикации сообщения
	log.Println("Создаем контекст с таймаутом 5 секунд...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Публикация сообщения в очередь
	// Параметры: exchange, routingKey, mandatory, immediate, msg
	log.Println("Публикуем сообщение в очередь...")
	message := "hello jack"
	err = ch.PublishWithContext(
		ctx,    // контекст с таймаутом
		"",     // exchange — пустая строка означает default exchange
		q.Name, // routingKey — имя очереди (для default exchange)
		false,  // mandatory — не требовать подтверждения доставки
		false,  // immediate — не требовать немедленной доставки
		amqp.Publishing{
			ContentType:  "text/plain",    // тип контента
			Body:         []byte(message), // тело сообщения
			DeliveryMode: amqp.Persistent, // сообщение будет сохранено на диск
		},
	)
	if err != nil {
		log.Fatalf("Ошибка публикации сообщения: %v", err)
	}

	log.Println("Сообщение '", message, "' успешно отправлено в очередь!")
	log.Println("Демонстрация завершена. Проверьте очередь в веб-панели: http://localhost:15672")
}
