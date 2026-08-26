package service

import (
	"errors"
	"fmt"
	"rune/internal/auth"
	"rune/internal/seal"
	"rune/internal/storage"
	"time"
)

type TokenService struct {
	store  storage.Store
	sealer *seal.Manager
}

func NewTokenService(store storage.Store, sealer *seal.Manager) *TokenService {
	return &TokenService{
		store:  store,
		sealer: sealer,
	}
}

func (s *TokenService) Create(name string, namespace string) (string, auth.TokenRecord, error) {

	if s.sealer.IsSealed() {
		return "", auth.TokenRecord{}, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	namespace = normalizeNamespace(namespace)

	if !s.store.NamespaceExists(namespace) {
		return "", auth.TokenRecord{}, errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	token, hash := auth.GenerateToken()
	id := fmt.Sprintf("tkn_%d", time.Now().UnixNano())

	record := auth.TokenRecord{
		ID:        id,
		Hash:      hash,
		Name:      name,
		Namespace: namespace,
		CreatedAt: time.Now().Format(time.RFC3339),
		Revoked:   false,
	}

	err := s.store.SaveToken(record)
	return token, record, err
}

func (s *TokenService) Revoke(id string) error {
	if s.sealer.IsSealed() {
		return errors.New("[OPERATION DENIED]: Vault is Sealed")
	}
	if id == "root" {
		return errors.New("[OPERATION DENIED]: Cannot revoke root token")
	}
	return s.store.RevokeToken(id)
}

func (s *TokenService) List() ([]auth.TokenRecord, error) {
	if s.sealer.IsSealed() {
		return nil, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}
	return s.store.ListTokens()
}
