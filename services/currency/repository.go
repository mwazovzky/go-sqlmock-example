package currency

import (
	"go-sqlmock-example/services/database"
	"time"
)

type Repository interface {
	Query(timeFrom time.Time) (database.RowsIterator, error)
	Parse(rows database.RowScanner) (Currency, error)
}

type CurrencyRepository struct {
	db database.DB
}

func NewCurrencyRepository(db database.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) Query(timeFrom time.Time) (database.RowsIterator, error) {
	query := "SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?"
	return r.db.Query(query, timeFrom)
}

func (r *CurrencyRepository) Parse(rows database.RowScanner) (Currency, error) {
	var currency Currency
	err := rows.Scan(&currency.Type, &currency.Chain, &currency.ISO, &currency.CreatedAt)
	if err != nil {
		return Currency{}, err
	}
	return currency, nil
}
