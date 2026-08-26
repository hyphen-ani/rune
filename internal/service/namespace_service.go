package service

import (
	"errors"
	"rune/internal/storage"
	"strings"
)

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
func (s *NamespaceService) DeleteNamespace(namespace string) error {

	if namespace == "default" {
		return errors.New("[INVALID OPERATION] Cannot delete default namespace")
	}

	keys, err := s.store.ListKeys()
	if err != nil {
		return err
	}

	for _, k := range keys {
		if strings.HasPrefix(k, namespace+"/") {
			return errors.New("[INVALID OPERATION] Cannot Delete Namespace, Namespace Not Empty")
		}
	}

	return s.store.DeleteNamespace(namespace)
}
