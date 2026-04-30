package service

import (
	"errors"
	"rune/internal/crypto"
	"rune/internal/seal"
	"rune/internal/storage"
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

func (s *SecretService) Put(key, value string) error {

	if s.sealer.IsSealed() {
		return errors.New("vault is sealed")
	}

	k := s.sealer.GetKey()

	ciphertext, nonce, err := crypto.Encrypt(k, []byte(value))
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

	if s.sealer.IsSealed() {
		return "", errors.New("vault is sealed")
	}

	k := s.sealer.GetKey()

	record, err := s.store.Get(key)
	if err != nil {
		return "", err
	}

	plaintext, err := crypto.Decrypt(k, record.Ciphertext, record.Nonce)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
