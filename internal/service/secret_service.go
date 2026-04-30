package service

import (
	"rune/internal/crypto"
	"rune/internal/storage"
)

type SecretService struct {
	store storage.Store
	key   []byte
}

func NewSecretService(store storage.Store, key []byte) *SecretService {
	return &SecretService{
		store: store,
		key:   key,
	}
}

func (s *SecretService) Put(key, value string) error {

	ciphertext, nonce, err := crypto.Encrypt(s.key, []byte(value))
	if err != nil {
		return err
	}

	record := storage.SecretRecord{
		Ciphertext: ciphertext,
		Nonce:      nonce,
	}

	return s.store.Put(key, record)
}

func (s *SecretService) Get(key string) (string, error) {

	record, err := s.store.Get(key)
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.Decrypt(s.key, record.Ciphertext, record.Nonce)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
