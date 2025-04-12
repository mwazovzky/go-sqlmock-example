package processor

import (
	"database/sql"
	"errors"
	"go-sqlmock-example/services/currency"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProducer implements currency.Producer for testing
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Produce(currency currency.Currency) error {
	args := m.Called(currency)
	return args.Error(0)
}

func TestProcessor_Process_Success(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	createdAt := time.Date(2025, time.April, 12, 23, 40, 31, 0, time.UTC)
	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"}).
		AddRow("fiat", "ethereum", "USD", createdAt)

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	mockProducer := new(MockProducer)
	mockProducer.On("Produce", currency.Currency{
		Type: "fiat", ISO: "USD", Chain: sql.NullString{String: "ethereum", Valid: true}, CreatedAt: createdAt,
	}).Return(nil)

	repo := currency.NewCurrencyRepository(db)
	processor := New(repo, mockProducer)

	// Execute
	err = processor.Process(time.Now())

	// Verify
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockProducer.AssertExpectations(t)
}

func TestProcessor_Process_QueryError(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(errors.New("query error"))

	mockProducer := new(MockProducer)

	repo := currency.NewCurrencyRepository(db)
	processor := New(repo, mockProducer)

	// Execute
	err = processor.Process(time.Now())

	// Verify
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockProducer.AssertExpectations(t)
}

func TestProcessor_Process_ParseError(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"}).
		AddRow("invalid", nil, nil, nil)

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	mockProducer := new(MockProducer)

	repo := currency.NewCurrencyRepository(db)
	processor := New(repo, mockProducer)

	// Execute
	err = processor.Process(time.Now())

	// Verify
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockProducer.AssertExpectations(t)
}

func TestProcessor_Process_ProduceError(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	createdAt := time.Date(2025, time.April, 12, 23, 43, 46, 504277000, time.UTC)
	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"}).
		AddRow("fiat", "ethereum", "USD", createdAt)

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	mockProducer := new(MockProducer)
	mockProducer.On("Produce", currency.Currency{
		Type: "fiat", ISO: "USD", Chain: sql.NullString{String: "ethereum", Valid: true}, CreatedAt: createdAt,
	}).Return(errors.New("produce error"))

	repo := currency.NewCurrencyRepository(db)
	processor := New(repo, mockProducer)

	// Execute
	err = processor.Process(time.Now())

	// Verify
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockProducer.AssertExpectations(t)
}

func TestProcessor_Process_NoRows(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"type", "chain", "iso", "created_at"})

	mock.ExpectQuery("SELECT type, chain, iso, created_at FROM currencies WHERE created_at > ?").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(mockRows)

	mockProducer := new(MockProducer)

	repo := currency.NewCurrencyRepository(db)
	processor := New(repo, mockProducer)

	// Execute
	err = processor.Process(time.Now())

	// Verify
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
	mockProducer.AssertExpectations(t)
}
