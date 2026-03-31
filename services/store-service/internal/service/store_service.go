package service

import (
	"context"
	"errors"
	"jualan-online/services/store-service/internal/domain"
)

type StoreService interface {
	CreateStore(ctx context.Context, name, platformType string) (*domain.Store, error)
	GetStoreByID(ctx context.Context, id string) (*domain.Store, error)
	GetAllStores(ctx context.Context) ([]domain.Store, error)
	UpdateStore(ctx context.Context, id, name, platformType string, isActive bool) (*domain.Store, error)
	DeleteStore(ctx context.Context, id string) error
}

type storeService struct {
	repo domain.StoreRepository
}

func NewStoreService(repo domain.StoreRepository) StoreService {
	return &storeService{repo: repo}
}

var validPlatforms = map[string]bool{
	"SHOPEE":      true,
	"TIKTOK_SHOP": true,
	"TOKOPEDIA":   true,
	"BLIBLI":      true,
	"LAZADA":       true,
	"ALFAGIFT":    true,
}

func (s *storeService) CreateStore(ctx context.Context, name, platformType string) (*domain.Store, error) {
	if !validPlatforms[platformType] {
		return nil, errors.New("invalid platform type")
	}

	store := &domain.Store{
		Name:         name,
		PlatformType: platformType,
		IsActive:     true,
	}

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *storeService) GetStoreByID(ctx context.Context, id string) (*domain.Store, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *storeService) GetAllStores(ctx context.Context) ([]domain.Store, error) {
	return s.repo.GetAll(ctx)
}

func (s *storeService) UpdateStore(ctx context.Context, id, name, platformType string, isActive bool) (*domain.Store, error) {
	store, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if platformType != "" && !validPlatforms[platformType] {
		return nil, errors.New("invalid platform type")
	}

	if name != "" {
		store.Name = name
	}
	if platformType != "" {
		store.PlatformType = platformType
	}
	store.IsActive = isActive

	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *storeService) DeleteStore(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
