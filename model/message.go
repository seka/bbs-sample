package model

import (
	"github.com/seka/bbs-sample/database"
)

// Message ...
type Message struct {
	ID        int
	UserID    int
	UserName  string
	Message   string
	CreatedAt string
}

// MessageModel ...
type MessageModel struct {
	db database.Database
}

// NewMessageModel ....
func NewMessageModel(db database.Database) *MessageModel {
	return &MessageModel{
		db: db,
	}
}

// FindAll ...
func (m *MessageModel) FindAll() ([]*Message, error) {
	query := `
	SELECT m.message message, m.created_at created_at, u.name name
	FROM messages m
	INNER JOIN users u ON m.user_id = u.id
	`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	messages := []*Message{}
	for rows.Next() {
		m := &Message{}
		if err := rows.Scan(&m.Message, &m.CreatedAt, &m.UserName); err != nil {
			break
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

// Save ...
func (m *MessageModel) Save(msg *Message) error {
	query := `INSERT INTO messages(user_id, message, created_at) VALUES (?, ?, ?)`
	_, err := m.db.Execute(query, msg.UserID, msg.Message, msg.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
