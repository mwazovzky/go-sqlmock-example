package main

import (
	"go-sqlmock-example/services/config"
	"go-sqlmock-example/services/currency"
	"go-sqlmock-example/services/database"
	"go-sqlmock-example/services/producer"
	"log"
	"time"
)

func main() {
	cfg := config.Load()

	db := database.OpenDB(cfg)
	defer database.CloseDB(db)

	cr := currency.NewCurrencyRepository(db)

	p := producer.NewProducer()
	cp := currency.NewCurrencyProducer(p)

	timeFrom, err := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	if err != nil {
		log.Fatal("Failed to parse time:", err)
	}

	err = Run(timeFrom, cr, cp)
	if err != nil {
		log.Fatal("Failed to run:", err)
	}
}

func Run(timeFrom time.Time, cr *currency.CurrencyRepository, cp *currency.CurrencyProducer) error {
	rows, err := cr.Query(timeFrom)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		currency, err := cr.Parse(rows)
		if err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		cp.Produce(currency)
		if err != nil {
			log.Fatal("Failed to produce currency:", err)
		}
	}

	return nil
}
