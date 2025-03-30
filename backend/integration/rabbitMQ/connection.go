package rabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	VHost    string
}

func Connect(cfg RabbitMQConfig) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	log.Println("RabbitMQ connected!")
	return conn, nil
}

func MustConnect(cfg RabbitMQConfig) *amqp.Connection {
	for {
		conn, err := Connect(cfg)
		if err != nil {
			log.Printf("Failed to connect: %v. Retry in 5 sec...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		return conn
	}
}
