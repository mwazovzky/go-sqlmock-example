package currency

import "encoding/json"

type MessageBroker interface {
	Produce([]byte) error
}

type Producer interface {
	Produce(Currency) error
}

type CurrencyProducer struct {
	producer MessageBroker
}

func NewCurrencyProducer(producer MessageBroker) *CurrencyProducer {
	return &CurrencyProducer{producer: producer}
}

func (p *CurrencyProducer) Produce(currency Currency) error {
	payload, err := json.Marshal(currency)
	if err != nil {
		return err
	}

	return p.producer.Produce(payload)
}
