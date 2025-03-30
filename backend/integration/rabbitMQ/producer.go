package rabbitMQ

import (
	"log"
	"workly/domain"
	"workly/usecase"

	"github.com/streadway/amqp"
)

type rabbitProducer struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// NewRabbitProducer – создаёт нового Producer
func NewRabbitProducer(conn *amqp.Connection, exchange, routingKey string) (usecase.IProducer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitProducer{
		channel:    ch,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

// Publish – реализация domain.Producer
func (p *rabbitProducer) Publish(msg domain.Message) error {
	log.Printf("Publishing message: %#v\n", msg)

	err := p.channel.Publish(
		p.exchange,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(msg.Content),
		},
	)
	return err
}
