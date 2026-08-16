
func main() {
	producer := NewProducer(conn, monitor)

	msg := &models.Message{
		ID:       1,
		Content:  "Hello with priority",
		Priority: 5,
	}

	err = producer.PublishJSON("", cfg.RabbitMQ.Queue, msg)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Println("Message sent successfully!")
	}
}