package database

import (
	"database/sql"
	"go-sqlmock-example/services/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// RowScanner is an interface that abstracts the Scan method.
// It allows for mocking and testing of row-scanning operations without relying on sql.Rows directly.
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// RowsIterator represents an interface for iterating over rows in a result set.
// It combines the functionality of moving to the next row (Next), closing the iterator (Close),
// and scanning the current row (RowScanner).
// This abstraction enables mocking and testing of row iteration logic.
type RowsIterator interface {
	Next() bool
	Close() error
	RowScanner
}

// DB is an interface for database operations.
// It abstracts the Query method to return a RowsIterator instead of sql.Rows,
// enabling easier testing and mocking of database interactions.
// The Close method ensures proper cleanup of database connections.
type DB interface {
	Query(query string, args ...interface{}) (RowsIterator, error)
	Close() error
}

// SQLRowsWrapper wraps sql.Rows to implement our RowsIterator interface
type SQLRowsWrapper struct {
	*sql.Rows
}

// Ensure SQLRowsWrapper implements RowsIterator
var _ RowsIterator = &SQLRowsWrapper{}

// Next moves to the next row
func (w *SQLRowsWrapper) Next() bool {
	return w.Rows.Next()
}

// Close closes the rows
func (w *SQLRowsWrapper) Close() error {
	return w.Rows.Close()
}

// SQLDBWrapper wraps sql.DB to implement the DB interface
type SQLDBWrapper struct {
	db *sql.DB
}

// Ensure SQLDBWrapper implements DB to guarantee compile-time interface compliance
var _ DB = &SQLDBWrapper{}

// Query executes a query and returns a RowsIterator
func (w *SQLDBWrapper) Query(query string, args ...interface{}) (RowsIterator, error) {
	rows, err := w.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLRowsWrapper{Rows: rows}, nil
}

// Close closes the database connection
func (w *SQLDBWrapper) Close() error {
	return w.db.Close()
}

// OpenDB opens a database connection and returns a DB interface
func OpenDB(cfg config.Config) DB {
	dsn := cfg.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	log.Printf("Opened DB connection: %s/%s", cfg.Host, cfg.Database)

	return &SQLDBWrapper{db: db}
}

// CloseDB closes a database connection
func CloseDB(db DB) {
	if err := db.Close(); err != nil {
		log.Fatal("Failed to close DB connection:", err)
	}

	log.Println("Closed DB connection")
}
