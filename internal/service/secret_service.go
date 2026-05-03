package service

import (
	"errors"
	"rune/internal/constants"
	"rune/internal/crypto"
	"rune/internal/seal"
	"rune/internal/storage"
	"strings"
)

type SecretService struct {
	store  storage.Store
	sealer *seal.Manager
}

func NewSecretService(store storage.Store, sealer *seal.Manager) *SecretService {
	return &SecretService{
		store:  store,
		sealer: sealer,
	}
}

func normalizeNamespace(ns string) string {
	if ns == "" {
		return constants.DefaultNamespace
	}
	return ns
}

func sanitizeKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.TrimPrefix(key, "/")
	key = strings.TrimPrefix(key, constants.DefaultNamespace+"/")
	return key
}

func (s *SecretService) Put(namespace, key, value string) error {

	if s.sealer.IsSealed() {
		return errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespace)

	if !s.store.NamespaceExists(ns) {
		return errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	cleanKey := sanitizeKey(key)

	if cleanKey == "" {
		return errors.New("[INVALID]: Invalid Key")
	}

	fullKey := ns + "/" + cleanKey

	k := s.sealer.GetKey()

	ciphertext, nonce, err := crypto.Encrypt(k, []byte(value))
	if err != nil {
		return err
	}

	record := storage.SecretRecord{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	}

	return s.store.Put(fullKey, record)
}

func (s *SecretService) Get(namespace, key string) (string, error) {

	if s.sealer.IsSealed() {
		return "", errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespace)

	cleanKey := sanitizeKey(key)

	if cleanKey == "" {
		return "", errors.New("[INVALID]: Invalid Key")
	}

	fullKey := ns + "/" + cleanKey

	k := s.sealer.GetKey()

	record, err := s.store.Get(fullKey)
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.Decrypt(k, record.Ciphertext, record.Nonce)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (s *SecretService) ListKeys(namespace string) ([]string, error) {
	if s.sealer.IsSealed() {
		return nil, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespace)

	if !s.store.NamespaceExists(ns) {
		return nil, errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	keys, err := s.store.ListKeys()
	if err != nil {
		return nil, err
	}

	var filtered []string

	for _, key := range keys {
		if strings.HasPrefix(key, ns+"/") {
			filtered = append(filtered, strings.TrimPrefix(key, ns+"/"))
		}
	}

	return filtered, nil
}
