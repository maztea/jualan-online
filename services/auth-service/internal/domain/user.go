package domain

import (
	"context"
	"time"
)

type User struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username     string    `json:"username" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         string    `json:"role" gorm:"not null;default:STAFF"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}
