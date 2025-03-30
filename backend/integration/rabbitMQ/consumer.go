package rabbitMQ

import (
	"github.com/streadway/amqp"
	"log"
)

// RabbitConsumer – отвечает за приём сообщений из очереди/обработку
type RabbitConsumer struct {
	channel   *amqp.Channel
	queueName string
}

// NewRabbitConsumer – создаёт Consumer и подписывается на очередь
func NewRabbitConsumer(conn *amqp.Connection, queueName string) (*RabbitConsumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Объявляем очередь
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitConsumer{
		channel:   ch,
		queueName: queueName,
	}, nil
}

// StartConsuming – запускает цикл чтения из очереди и обработки
func (c *RabbitConsumer) StartConsuming() {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	go func() {
		for msg := range msgs {
			log.Printf("Received message: %s", msg.Body)
		}
	}()
}
