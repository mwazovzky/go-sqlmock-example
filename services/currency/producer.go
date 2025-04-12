package currency

import "encoding/json"

type Producer interface {
	Produce([]byte) error
}

type CurrencyProducer struct {
	producer Producer
}

func NewCurrencyProducer(producer Producer) *CurrencyProducer {
	return &CurrencyProducer{producer: producer}
}

func (p *CurrencyProducer) Produce(currency Currency) error {
	payload, err := json.Marshal(currency)
	if err != nil {
		return err
	}

	return p.producer.Produce(payload)
}
