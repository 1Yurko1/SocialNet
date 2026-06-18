package models

import "github.com/google/uuid"

// TokenClaims используется для JWT
type TokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
}

// UserSession для хранения в Redis
type UserSession struct {
	UserID        uuid.UUID `json:"user_id"`
	LastSeen      int64     `json:"last_seen"`
	CurrentChatID uuid.UUID `json:"current_chat_id,omitempty"`
}

// WSMessage — общая структура для WebSocket сообщений
type WSMessage struct {
	Type      string `json:"type"` // "private_msg", "group_msg", "status"
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ToUserID  string `json:"to_user_id,omitempty"` // Для приватных сообщений
	SenderID  string `json:"sender_id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
