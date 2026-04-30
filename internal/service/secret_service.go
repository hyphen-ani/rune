package service

import "errors"

type SecretService struct {
	store map[string]string
}

func NewSecretService() *SecretService {
	return &SecretService{
		store: make(map[string]string),
	}
}

func (s *SecretService) Put(key, value string) {
	s.store[key] = value
}

func (s *SecretService) Get(key string) (string, error) {
	val, ok := s.store[key]
	if !ok {
		return "", errors.New("secret not found")
	}
	return val, nil
}
