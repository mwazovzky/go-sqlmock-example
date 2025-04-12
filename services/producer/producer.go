package producer

import "log"

type Producer struct{}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Produce(data []byte) error {
	log.Println("Producing data:", string(data))
	return nil
}
