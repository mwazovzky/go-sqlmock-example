package main

import (
	"database/sql"
	"go-sqlmock-example/services/config"
	"go-sqlmock-example/services/database"
	"log"
	"time"
)

type Currency struct {
	Type      string
	ISO       string
	Chain     sql.NullString
	CreatedAt time.Time
}

func main() {
	cfg := config.Load()

	db := database.OpenDB(cfg)
	defer database.CloseDB(db)

	timeFrom, err := time.Parse("2006-01-02 15:04:05", "2025-01-01 00:00:00")
	if err != nil {
		log.Fatal("Failed to parse time:", err)
	}

	query := "SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?"
	rows, err := db.Query(query, timeFrom)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var currency Currency

		if err := rows.Scan(&currency.Type, &currency.Chain, &currency.ISO, &currency.CreatedAt); err != nil {
			log.Fatal("Failed to scan row:", err)
		}

		log.Printf("currency: %+v, ", currency)
	}
}
