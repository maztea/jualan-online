package repository

import (
	"context"
	"gorm.io/gorm"
	"jualan-online/services/store-service/internal/domain"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) domain.StoreRepository {
	return &storeRepository{db: db}
}

func (r *storeRepository) Create(ctx context.Context, store *domain.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}

func (r *storeRepository) GetByID(ctx context.Context, id string) (*domain.Store, error) {
	var store domain.Store
	err := r.db.WithContext(ctx).First(&store, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storeRepository) GetAll(ctx context.Context) ([]domain.Store, error) {
	var stores []domain.Store
	err := r.db.WithContext(ctx).Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *storeRepository) Update(ctx context.Context, store *domain.Store) error {
	return r.db.WithContext(ctx).Save(store).Error
}

func (r *storeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Store{}, "id = ?", id).Error
}
