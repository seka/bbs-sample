package database

import (
	"database/sql"
	"errors"
)

var (
	// ErrConnNotExist ...
	ErrConnNotExist = errors.New("database: connection is not exists")
)

// Database ...
type Database interface {
	// Connnect ...
	Connect() error

	// Query ...
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// Execute ...
	Execute(query string, args ...interface{}) (sql.Result, error)

	// Disconnect ...
	Disconnect() error

	// Ping ...
	Ping() error
}
