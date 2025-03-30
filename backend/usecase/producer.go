package usecase

import "workly/domain"

type IProducer interface {
	Publish(msg domain.Message) error
}

type ProducerUsecase struct {
	producer IProducer
}

func NewProducerUsecase(prod IProducer) *ProducerUsecase {
	return &ProducerUsecase{
		producer: prod,
	}
}

func (uc *ProducerUsecase) SendMessage(message domain.Message) error {
	return uc.producer.Publish(message)
}
