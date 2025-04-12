package mocks

import (
	"go-sqlmock-example/services/database"

	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of the database.DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(query string, args ...interface{}) (database.RowsIterator, error) {
	arguments := m.Called(query, args)
	if arguments.Get(0) == nil {
		return nil, arguments.Error(1)
	}
	return arguments.Get(0).(database.RowsIterator), arguments.Error(1)
}

func (m *MockDB) Close() error {
	args := m.Called()
	return args.Error(0)
}
