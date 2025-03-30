package app

import (
	"log"
	"time"
	"workly/domain"
	"workly/integration/rabbitMQ"
	"workly/usecase"
)

func NewRabbitMQ(rabbitCfg rabbitMQ.RabbitMQConfig) {
	conn := rabbitMQ.MustConnect(rabbitCfg)
	defer conn.Close()

	producer, err := rabbitMQ.NewRabbitProducer(conn, "test-exchange", "test-routing")
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	consumer, err := rabbitMQ.NewRabbitConsumer(conn, "test-queue")
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	consumer.StartConsuming()

	producerUC := usecase.NewProducerUsecase(producer)

	for i := 1; i <= 5; i++ {
		message := domain.Message{
			ID:      "msg_id_123",
			Content: "Hello from Clean Architecture!",
		}
		if err := producerUC.SendMessage(message); err != nil {
			log.Printf("Failed to send message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

}
