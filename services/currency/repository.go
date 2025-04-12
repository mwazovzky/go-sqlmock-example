package currency

import (
	"database/sql"
	"time"
)

type CurrencyRepository struct {
	db *sql.DB
}

func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) Query(timeFrom time.Time) (*sql.Rows, error) {
	query := "SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?"
	return r.db.Query(query, timeFrom)
}

func (r *CurrencyRepository) Parse(rows *sql.Rows) (Currency, error) {
	var currency Currency
	err := rows.Scan(&currency.Type, &currency.Chain, &currency.ISO, &currency.CreatedAt)
	if err != nil {
		return Currency{}, err
	}
	return currency, nil
}
