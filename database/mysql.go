package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // for sql.Open
)

// Options ...
type Options struct {
	Addr     string
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
		dsn: fmt.Sprint(opt.User, ":", opt.Password, "@tcp(", opt.Addr, ")/", opt.Name),
	}
}

// Connect ...
func (m *MySQL) Connect() error {
	conn, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return err
	}
	m.conn = conn
	return m.Ping()
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

// Ping ...
func (m *MySQL) Ping() error {
	if m.conn == nil {
		return ErrConnNotExist
	}
	return m.conn.Ping()
}

var _ Database = (*MySQL)(nil)
