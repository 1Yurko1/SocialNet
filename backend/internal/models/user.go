package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"` // Уникальный ник
	Email     string    `json:"email"`    // Уникальный email
	Password  string    `json:"-"`        // Хэш пароля, никогда не возвращаем в JSON
	AvatarURL string    `json:"avatar_url"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Relation struct {
	FollowerID  uuid.UUID `json:"follower_id"`  // Кто подписывается
	FollowingID uuid.UUID `json:"following_id"` // На кого подписываются
	CreatedAt   time.Time `json:"created_at"`
}
