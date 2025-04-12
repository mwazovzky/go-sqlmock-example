package currency

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyRepository_Query(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"}).
		AddRow("fiat", "ethereum", "USD", time.Now())

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	repo := NewCurrencyRepository(db)
	timeFrom := time.Now()
	rows, err := repo.Query(timeFrom)

	assert.NoError(t, err)
	assert.NotNil(t, rows)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurrencyRepository_Query_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	repo := NewCurrencyRepository(db)
	timeFrom := time.Now()
	_, err = repo.Query(timeFrom)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurrencyRepository_Parse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"}).
		AddRow("fiat", "ethereum", "USD", time.Now())

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	repo := NewCurrencyRepository(db)

	// Use the mock rows to test Parse
	rows, err := repo.Query(time.Now())
	assert.NoError(t, err)

	for rows.Next() {
		currency, err := repo.Parse(rows)
		assert.NoError(t, err)
		assert.Equal(t, "fiat", currency.Type)
		assert.Equal(t, "ethereum", currency.Chain.String)
		assert.Equal(t, "USD", currency.ISO)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCurrencyRepository_Parse_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"}).
		AddRow("invalid", nil, nil, nil)

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	repo := NewCurrencyRepository(db)

	// Use the mock rows to test Parse
	rows, err := repo.Query(time.Now())
	assert.NoError(t, err)

	for rows.Next() {
		_, err := repo.Parse(rows)
		assert.Error(t, err)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}
