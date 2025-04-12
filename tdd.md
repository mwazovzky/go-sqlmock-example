# Testing Database Operations in Go

This document explains the approach used in this project to implement and test database operations. It focuses on the use of interfaces and mocks to achieve testability.

## Approach

### 1. Abstraction with Interfaces

To make database operations testable, we abstracted key functionalities using interfaces:

- **`DB` Interface**: Represents the database connection and provides a `Query` method that returns a `RowsIterator`.
- **`RowsIterator` Interface**: Combines row iteration (`Next`), cleanup (`Close`), and row scanning (`RowScanner`).
- **`RowScanner` Interface**: Abstracts the `Scan` method for reading row data.

These abstractions decouple the application logic from the `database/sql` package, enabling us to mock database behavior in tests.

### 2. Wrappers for `database/sql`

Since `database/sql` does not natively implement our interfaces, we created lightweight wrappers:

- **`SQLDBWrapper`**: Wraps `*sql.DB` to implement the `DB` interface.
- **`SQLRowsWrapper`**: Wraps `*sql.Rows` to implement the `RowsIterator` interface.

These wrappers bridge the gap between the `database/sql` package and our custom interfaces.

### 3. Mocking for Unit Tests

To test the application logic without relying on a real database, we implemented mocks:

- **`MockDB`**: Simulates the `DB` interface.
- **`MockRows`**: Simulates the `RowsIterator` interface.
- **`MockRowScanner`**: Simulates the `RowScanner` interface.

These mocks allow us to test various scenarios, such as:

- Successful queries and row processing.
- Errors during queries or row scanning.
- Edge cases like empty result sets.

### 4. Trade-offs of This Approach

While this approach provides excellent testability, it has some potential disadvantages:

- **Additional Complexity**: Introducing custom interfaces and wrappers adds complexity to the codebase.
- **Limited `database/sql` Features**: If the application needs to use many features of `database/sql` (e.g., transactions, prepared statements), we would need to extend our interfaces and mocks, increasing maintenance overhead.
- **Mock Maintenance**: For every new database operation, corresponding mocks must be implemented, which can be time-consuming.

## Example: Testing the Processor

The `CurrencyProcessor` demonstrates how to use the abstractions and mocks:

1. The `Processor` queries the database using the `DB` interface.
2. It iterates over rows using the `RowsIterator` interface.
3. It scans row data using the `RowScanner` interface.

In tests, we replace the real database with mocks to simulate various scenarios.

### Test Example

```go
func TestProcessor_Process_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockRows := mocks.NewMockRows([]bool{true, false}) // One row followed by end
	currencyObj := currency.Currency{Type: "fiat", ISO: "USD"}

	mockRepo.On("Query", mock.Anything).Return(mockRows, nil)
	mockRows.On("Close").Return(nil)
	mockRepo.On("Parse", mockRows).Return(currencyObj, nil)
	mockProducer.On("Produce", currencyObj).Return(nil)

	processor := New(mockRepo, mockProducer)
	err := processor.Process(time.Now())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}
```

## Conclusion

This project demonstrates how to write testable code for database operations in Go. By abstracting database interactions and using mocks, we can achieve high test coverage and ensure the reliability of our application logic. However, developers should weigh the trade-offs of this approach and consider the complexity it introduces.
