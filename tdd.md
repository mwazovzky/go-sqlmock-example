# Testing Database Operations in Go with `sqlmock`

This document explains how to test database operations in Go using the `sqlmock` library. By directly using `database/sql` and `sqlmock`, we simplify the codebase while maintaining testability.

## Approach

### 1. Direct Use of `database/sql`

Instead of creating custom interfaces and wrappers, we directly use the `database/sql` package for database operations. This reduces complexity and ensures compatibility with the full feature set of `database/sql`.

### 2. Mocking with `sqlmock`

The `sqlmock` library allows us to mock `*sql.DB` and `*sql.Rows` for testing. This eliminates the need for custom mocks and abstractions.

### 3. Writing Tests with `sqlmock`

With `sqlmock`, we can simulate various database behaviors, such as:

- Successful queries and row processing.
- Errors during queries or row scanning.
- Edge cases like empty result sets.

### Example: Testing the Repository

The `CurrencyRepository` demonstrates how to use `sqlmock` for testing database operations.

#### Test Example

```go
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
```

### Trade-offs of This Approach

While this approach simplifies the codebase, it has some potential disadvantages:

- **Dependency on `sqlmock`**: The tests rely on the `sqlmock` library, which may require updates if the library changes.
- **Limited Abstraction**: Without custom interfaces, the code is tightly coupled to `database/sql`, which may reduce flexibility in some scenarios.

## Conclusion

Using `sqlmock` simplifies the codebase while maintaining testability. This approach is ideal for projects where simplicity and compatibility with `database/sql` are priorities.
