package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatType string

const (
	ChatPrivate ChatType = "private" // Личка
	ChatGroup   ChatType = "group"   // Группа
)

type Chat struct {
	ID        uuid.UUID `json:"id"`
	Type      ChatType  `json:"type"`
	Title     string    `json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMember struct {
	ChatID   uuid.UUID `json:"chat_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type Message struct {
	ID          uuid.UUID `json:"id"`
	ChatID      uuid.UUID `json:"chat_id"`
	SenderID    uuid.UUID `json:"sender_id"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}
