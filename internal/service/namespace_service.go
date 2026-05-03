package service

import "rune/internal/storage"

type NamespaceService struct {
	store storage.Store
}

func NewNamespaceService(store storage.Store) *NamespaceService {
	return &NamespaceService{
		store: store,
	}
}

func (s *NamespaceService) CreateNamespace(namespace string) error {
	return s.store.CreateNamespace(namespace)
}

func (s *NamespaceService) ListNamespace() ([]string, error) {
	return s.store.ListNamespaces()
}
