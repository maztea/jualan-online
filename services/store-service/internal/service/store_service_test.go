package service

import (
	"context"
	"errors"
	"jualan-online/services/store-service/internal/domain"
	"testing"
)

type mockRepo struct {
	domain.StoreRepository
	stores map[string]*domain.Store
}

func (m *mockRepo) Create(ctx context.Context, store *domain.Store) error {
	m.stores[store.Name] = store
	return nil
}

func (m *mockRepo) GetByID(ctx context.Context, id string) (*domain.Store, error) {
	for _, s := range m.stores {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func TestStoreService_CreateStore(t *testing.T) {
	repo := &mockRepo{stores: make(map[string]*domain.Store)}
	svc := NewStoreService(repo)

	tests := []struct {
		name         string
		platformType string
		wantErr      bool
	}{
		{"Valid Store", "SHOPEE", false},
		{"Valid Store 2", "TOKOPEDIA", false},
		{"Invalid Store", "AMAZON", true},
		{"Empty Platform", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateStore(context.Background(), tt.name, tt.platformType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
