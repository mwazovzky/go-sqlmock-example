package main

import (
	"go-sqlmock-example/services/config"
	"go-sqlmock-example/services/currency"
	"go-sqlmock-example/services/database"
	"go-sqlmock-example/services/processor"
	"go-sqlmock-example/services/producer"
	"log"
	"time"
)

func main() {
	cfg := config.Load()

	db := database.OpenDB(cfg)
	defer database.CloseDB(db)

	currencyRepository := currency.NewCurrencyRepository(db)

	p := producer.NewProducer()
	currencyProducer := currency.NewCurrencyProducer(p)

	timeFrom, err := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	if err != nil {
		log.Fatal("Failed to parse time:", err)
	}

	proc := processor.New(currencyRepository, currencyProducer)
	err = proc.Process(timeFrom)
	if err != nil {
		log.Fatalf("Failed to process currencies: %v", err)
	}
}
