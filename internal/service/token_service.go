package service

import (
	"fmt"
	"rune/internal/auth"
	"rune/internal/storage"
	"time"
)

type TokenService struct {
	store storage.Store
}

func NewTokenService(store storage.Store) *TokenService {
	return &TokenService{
		store: store,
	}
}

func (s *TokenService) Create(name string) (string, auth.TokenRecord, error) {
	token, hash := auth.GenerateToken()
	id := fmt.Sprintf("tkn_%d", time.Now().UnixNano())

	record := auth.TokenRecord{
		ID:        id,
		Hash:      hash,
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
		Revoked:   false,
	}

	err := s.store.SaveToken(record)
	return token, record, err
}

func (s *TokenService) Revoke(id string) error {
	return s.store.RevokeToken(id)
}

func (s *TokenService) List() ([]auth.TokenRecord, error) {
	return s.store.ListTokens()
}
