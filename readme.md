# Go SQLMock Example

This project demonstrates how to test database operations in Go using interfaces and mocks.

## Key Features

- **Abstraction with Interfaces**: The `DB`, `RowsIterator`, and `RowScanner` interfaces abstract database operations, making the code testable and decoupled from `database/sql`.
- **Mocking for Unit Tests**: Custom mocks are implemented to simulate database behavior, enabling comprehensive unit testing.
- **Separation of Concerns**: The project follows a modular architecture, separating database logic, business logic, and testing concerns.

## How to Run

1. Set up the environment variables for MySQL as described in the `docker-compose.yaml`.
2. Start the MySQL database using Docker:
   ```bash
   docker-compose up
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
4. Run the tests:
   ```bash
   go test ./...
   ```

## Key Objectives

- Abstract database operations using interfaces.
- Use mocks to test database interactions.
- Explore the trade-offs of using custom abstractions for testing.

## Config

```bash
export MYSQL_HOST=localhost
export MYSQL_TCP_PORT=3306
export MYSQL_DATABASE=example
export MYSQL_ROOT_PASSWORD=rootsecret
export MYSQL_USER=user
export MYSQL_PASSWORD=usersecret
```
