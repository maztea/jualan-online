package domain

import (
	"context"
	"time"
)

type Store struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name         string    `json:"name" gorm:"not null"`
	PlatformType string    `json:"platform_type" gorm:"not null"` // SHOPEE, TOKOPEDIA, etc
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StoreRepository interface {
	Create(ctx context.Context, store *Store) error
	GetByID(ctx context.Context, id string) (*Store, error)
	GetAll(ctx context.Context) ([]Store, error)
	Update(ctx context.Context, store *Store) error
	Delete(ctx context.Context, id string) error
}
