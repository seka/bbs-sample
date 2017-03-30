package database

import (
	"context"
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
	Connect(ctx context.Context) error

	// Query ...
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// Execute ...
	Execute(query string, args ...interface{}) (sql.Result, error)

	// Disconnect ...
	Disconnect() error
}
