# Go SQLMock Example

This project demonstrates how to test database operations in Go using the `sqlmock` library. By directly using `database/sql` and `sqlmock`, we simplify the codebase while maintaining testability.

## Key Features

- **Direct Use of `database/sql`**: The project uses `database/sql` without custom wrappers or abstractions.
- **Mocking with `sqlmock`**: The `sqlmock` library is used to simulate database behavior, eliminating the need for a custom `mocks` package.
- **Simplified Codebase**: By removing custom interfaces, wrappers, and mocks, the code is easier to understand and maintain.

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

- Simplify database operations by directly using `database/sql`.
- Use `sqlmock` to test database interactions.
- Maintain test coverage while reducing code complexity.

## Config

```bash
export MYSQL_HOST=localhost
export MYSQL_TCP_PORT=3306
export MYSQL_DATABASE=example
export MYSQL_ROOT_PASSWORD=rootsecret
export MYSQL_USER=user
export MYSQL_PASSWORD=usersecret
```
