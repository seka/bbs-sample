package model

import (
	"github.com/seka/bbs-sample/database"
)

// User ...
type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

// UserModel ...
type UserModel struct {
	db database.Database
}

// NewUserModel ...
func NewUserModel(db database.Database) *UserModel {
	return &UserModel{
		db: db,
	}
}

// Find ...
func (u *UserModel) Find(user *User) (*User, error) {
	query := `
	SELECT id, name, email FROM users
	WHERE email=?
	AND password_hash=?
	LIMIT 1
	`
	rows, err := u.db.Query(query, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return user, nil
}

// FindAll ...
func (u *UserModel) FindAll() ([]*User, error) {
	query := `SELECT id, name, email FROM users`
	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}
	users := []*User{}
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			break
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// Save ...
func (u *UserModel) Save(user *User) error {
	query := `INSERT INTO users(id, name, email, password_hash) VALUES (?, ?, ?, ?)`
	_, err := u.db.Execute(query, user.ID, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// Exists ...
func (u *UserModel) Exists(user *User) bool {
	query := `SELECT TRUE  users(id) VALUES (?, ?)`
	rows, err := u.db.Query(query, user.Email, user.Password)
	if err != nil {
		return false
	}
	return rows.Next()
}
