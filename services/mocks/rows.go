package mocks

import (
	"go-sqlmock-example/services/database"

	"github.com/stretchr/testify/mock"
)

// MockRowScanner mocks the RowScanner interface for testing
type MockRowScanner struct {
	mock.Mock
}

func (m *MockRowScanner) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

// MockRows implements database.RowsIterator for testing
// It simulates a database cursor that can iterate through rows
type MockRows struct {
	mock.Mock
	currentIndex int
	hasNext      []bool // Controls the behavior of Next() method
}

// Ensure MockRows implements database.RowsIterator
var _ database.RowsIterator = &MockRows{}

// NewMockRows creates a new MockRows with predefined row availability
func NewMockRows(hasNext []bool) *MockRows {
	return &MockRows{
		currentIndex: -1,
		hasNext:      hasNext,
	}
}

// Next simulates cursor movement to the next row
func (m *MockRows) Next() bool {
	m.currentIndex++
	if m.currentIndex < len(m.hasNext) {
		return m.hasNext[m.currentIndex]
	}
	return false
}

// Close simulates closing the result set
func (m *MockRows) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Scan simulates reading column values into the provided destination variables
func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}
