package processor

import (
	"go-sqlmock-example/services/currency"
	"time"
)

// CurrencyProcessor handles the retrieval and processing of currencies.
// It connects the repository (for data retrieval) and producer (for data handling).
type CurrencyProcessor struct {
	repository currency.Repository
	producer   currency.Producer
}

// New creates a new CurrencyProcessor with dependencies.
func New(repository currency.Repository, producer currency.Producer) *CurrencyProcessor {
	return &CurrencyProcessor{
		repository: repository,
		producer:   producer,
	}
}

// Process retrieves and processes currencies created after a specific time.
// It queries the database, parses each row into a Currency object,
// and sends each currency to the producer.
func (p *CurrencyProcessor) Process(timeFrom time.Time) error {
	rows, err := p.repository.Query(timeFrom)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		currency, err := p.repository.Parse(rows)
		if err != nil {
			return err
		}

		if err = p.producer.Produce(currency); err != nil {
			return err
		}
	}

	return rows.Err() // Check for errors encountered during iteration
}
