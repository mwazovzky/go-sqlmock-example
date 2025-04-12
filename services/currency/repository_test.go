package currency

import (
	"database/sql"
	"errors"
	"go-sqlmock-example/services/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCurrencyRepository_Query(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockRows := &sql.Rows{}
	mockDB.On("Query", "SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?", mock.Anything).Return(mockRows, nil)

	repo := NewCurrencyRepository(mockDB)
	timeFrom := time.Now()
	rows, err := repo.Query(timeFrom)

	assert.NoError(t, err)
	assert.Equal(t, mockRows, rows)
	mockDB.AssertExpectations(t)
}

func TestCurrencyRepository_Query_Error(t *testing.T) {
	mockDB := new(mocks.MockDB)
	mockError := errors.New("query error")
	mockDB.On("Query", "SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?", mock.Anything).Return(nil, mockError)

	repo := NewCurrencyRepository(mockDB)
	timeFrom := time.Now()
	_, err := repo.Query(timeFrom)

	assert.Error(t, err)
	assert.Equal(t, mockError, err)
	mockDB.AssertExpectations(t)
}

func TestCurrencyRepository_Parse(t *testing.T) {
	mockDB := new(mocks.MockDB)
	repo := NewCurrencyRepository(mockDB)

	currencyType := "fiat"
	currencyChain := sql.NullString{String: "ethereum", Valid: true}
	currencyISO := "EUR"
	createdAt := time.Now()

	mockRowScanner := new(mocks.MockRowScanner)
	mockRowScanner.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		*args.Get(0).(*string) = currencyType
		*args.Get(1).(*sql.NullString) = currencyChain
		*args.Get(2).(*string) = currencyISO
		*args.Get(3).(*time.Time) = createdAt
	})

	currency, err := repo.Parse(mockRowScanner)

	assert.NoError(t, err)
	assert.Equal(t, currencyType, currency.Type)
	assert.Equal(t, currencyChain, currency.Chain)
	assert.Equal(t, currencyISO, currency.ISO)
	assert.Equal(t, createdAt, currency.CreatedAt)

	mockDB.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}

func TestCurrencyRepository_Parse_Error(t *testing.T) {
	mockDB := new(mocks.MockDB)
	repo := NewCurrencyRepository(mockDB)
	mockError := errors.New("scan error")

	mockRowScanner := new(mocks.MockRowScanner)
	mockRowScanner.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockError)

	_, err := repo.Parse(mockRowScanner)

	assert.Error(t, err)
	assert.Equal(t, mockError, err)

	mockDB.AssertExpectations(t)
	mockRowScanner.AssertExpectations(t)
}
