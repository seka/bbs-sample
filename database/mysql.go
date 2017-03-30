package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // for sql.Open
	"github.com/inconshreveable/log15"
)

// Options ...
type Options struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// MySQL ...
type MySQL struct {
	dsn  string
	conn *sql.DB
}

// NewMySQL ...
func NewMySQL(opt Options) Database {
	return &MySQL{
		dsn: fmt.Sprint(opt.User, ":", opt.Password, "@tcp(", opt.Host, ":", opt.Port, ")/", opt.Name),
	}
}

// Connect ...
func (m *MySQL) Connect(ctx context.Context) error {
	conn, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return err
	}
	m.conn = conn
	healthCheckErrCh := make(chan error, 1)
	go func() {
		healthCheckErrCh <- m.periodicHealthCheck(ctx)
	}()
	select {
	case err := <-healthCheckErrCh:
		return err
	case <-ctx.Done():
		<-healthCheckErrCh // Wait for finish
		if err := m.Disconnect(); err != nil {
			return err
		}
		return ctx.Err()
	}
}

// Query ...
func (m *MySQL) Query(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := m.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Execute ...
func (m *MySQL) Execute(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := m.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Disconnect ...
func (m *MySQL) Disconnect() error {
	if m.conn == nil {
		return ErrConnNotExist
	}
	if err := m.conn.Close(); err != nil {
		return err
	}
	m.conn = nil
	return nil
}

func (m *MySQL) periodicHealthCheck(ctx context.Context) error {
	if m.conn == nil {
		return ErrConnNotExist
	}
	if err := m.conn.Ping(); err != nil {
		return err
	}
	log15.Info("Connected to the database")
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := m.conn.Ping(); err != nil {
				return err
			}
		case <-ctx.Done():
			ticker.Stop()
			return ctx.Err()
		}
	}
}

var _ Database = (*MySQL)(nil)
