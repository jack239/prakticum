
func main() {
	consumer := NewConsumer(conn, monitor)

	handler := func(msg *models.Message) bool {
		fmt.Printf("Received: ID=%d, Content='%s', Priority=%d\n", msg.ID, msg.Content, msg.Priority)

		// Имитируем ошибку для части сообщений
		if msg.ID%2 == 0 {
			fmt.Println("  → Processing failed, sending to DLX")
			return false // NACK → сообщение уйдёт в DLX
		}

		fmt.Println("  → Processed successfully")
		return true // ACK
	}

	err = consumer.ConsumeJSON(cfg.RabbitMQ.Queue, handler)
	if err != nil {
		log.Fatal("Failed to consume messages: ", err)
	}
}