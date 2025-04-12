package processor

import (
	"errors"
	"go-sqlmock-example/services/currency"
	"go-sqlmock-example/services/database"
	"go-sqlmock-example/services/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository implements currency.Repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Query(timeFrom time.Time) (database.RowsIterator, error) {
	args := m.Called(timeFrom)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(database.RowsIterator), args.Error(1)
}

func (m *MockRepository) Parse(rows database.RowScanner) (currency.Currency, error) {
	args := m.Called(rows)
	return args.Get(0).(currency.Currency), args.Error(1)
}

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
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockRows := mocks.NewMockRows([]bool{true, false}) // One row followed by end
	currencyObj := currency.Currency{Type: "fiat", ISO: "USD"}

	mockRepo.On("Query", mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRepo.On("Parse", mockRows).Return(currencyObj, nil)
	mockProducer.On("Produce", currencyObj).Return(nil)

	// Execute
	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	// Verify
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestProcessor_Process_QueryError(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockError := errors.New("query error")

	mockRepo.On("Query", mock.Anything).Return(nil, mockError)

	// Execute
	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	// Verify
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

func TestProcessor_Process_ParseError(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockRows := mocks.NewMockRows([]bool{true})
	mockError := errors.New("parse error")

	mockRepo.On("Query", mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRepo.On("Parse", mockRows).Return(currency.Currency{}, mockError)

	// Execute
	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	// Verify
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestProcessor_Process_ProduceError(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockRows := mocks.NewMockRows([]bool{true})
	currencyObj := currency.Currency{Type: "fiat", ISO: "USD"}
	mockError := errors.New("produce error")

	mockRepo.On("Query", mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRepo.On("Parse", mockRows).Return(currencyObj, nil)
	mockProducer.On("Produce", currencyObj).Return(mockError)

	// Execute
	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	// Verify
	assert.Error(t, err)
	assert.Equal(t, mockError, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestProcessor_Process_NoRows(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockRows := mocks.NewMockRows([]bool{false}) // No rows

	mockRepo.On("Query", mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)

	// Execute
	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	// Verify
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

func TestProcessor_Process_MultipleRows(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockRows := mocks.NewMockRows([]bool{true, true, false}) // Two rows followed by end

	currency1 := currency.Currency{Type: "fiat", ISO: "USD"}
	currency2 := currency.Currency{Type: "crypto", ISO: "BTC"}

	mockRepo.On("Query", mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)

	mockRepo.On("Parse", mockRows).Return(currency1, nil).Once()
	mockProducer.On("Produce", currency1).Return(nil).Once()

	mockRepo.On("Parse", mockRows).Return(currency2, nil).Once()
	mockProducer.On("Produce", currency2).Return(nil).Once()

	// Execute
	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	// Verify
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}
