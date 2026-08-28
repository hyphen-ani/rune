package service

import (
	"errors"
	"log"
	"rune/internal/auth"
	"rune/internal/constants"
	"rune/internal/crypto"
	"rune/internal/seal"
	"rune/internal/storage"
	"strings"
	"time"
)

type SecretService struct {
	store  storage.Store
	sealer *seal.Manager
}

type SecretVersion struct {
	Version   int    `json:"version"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	RotatedAt string `json:"rotated_at"`
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

	log.Printf(

		"PUT DEBUG -> namespace=%q key=%q fullKey=%q",
		ns,
		key,
		fullKey,
	)
	now := time.Now().Format(time.RFC3339)

	record := storage.SecretRecord{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		CreatedAt:  now,
		UpdatedAt:  now,
		RotatedAt:  "",
		Version:    1,
	}

	versionKey := storage.VersionedSecretKey(fullKey, 1)

	if err := s.store.Put(versionKey, record); err != nil {
		return err
	}

	if err := s.store.Put(fullKey, record); err != nil {
		return err
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
func (s *SecretService) Delete(namespace, key string) error {
	if s.sealer.IsSealed() {
		return errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespace)

	if !s.store.NamespaceExists(ns) {
		return errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	fullKey := ns + "/" + key

	log.Printf(
		"DELETE DEBUG -> namespace=%q key=%q fullKey=%q",
		ns,
		key,
		fullKey,
	)

	current, err := s.store.Get(fullKey)
	if err != nil {
		return err
	}

	for version := 1; version <= current.Version; version++ {
		versionKey := storage.VersionedSecretKey(fullKey, version)
		_ = s.store.Delete(versionKey)
	}

	if err := s.store.Delete(fullKey); err != nil {
		return err
	}

	return nil
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
			if storage.IsVersionedSecretKey(key) {
				continue
			}
			filtered = append(filtered, strings.TrimPrefix(key, ns+"/"))
		}
	}

	return filtered, nil
}

func (s *SecretService) Rotate(namespace string, key string) (string, error) {
	if s.sealer.IsSealed() {
		return "", errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespace)

	if !s.store.NamespaceExists(ns) {
		return "", errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	fullKey := ns + "/" + key

	existing, err := s.store.Get(fullKey)
	if err != nil {
		return "", err
	}

	newValue, err := auth.GenerateSecret()
	if err != nil {
		return "", err
	}

	k := s.sealer.GetKey()
	ciphertext, nonce, err := crypto.Encrypt(k, []byte(newValue))
	if err != nil {
		return "", err
	}

	now := time.Now().Format(time.RFC3339)
	newVersion := existing.Version + 1
	record := storage.SecretRecord{
		Ciphertext: ciphertext,
		Nonce:      nonce,
		CreatedAt:  existing.CreatedAt,
		UpdatedAt:  now,
		RotatedAt:  now,
		Version:    newVersion,
	}

	versionKey := storage.VersionedSecretKey(fullKey, newVersion)

	err = s.store.Put(versionKey, record)
	if err != nil {
		return "", err
	}

	err = s.store.Put(fullKey, record)
	if err != nil {
		return "", err
	}

	return newValue, nil
}
func (s *SecretService) GetVersion(namespace, key string, version int) (string, error) {
	if s.sealer.IsSealed() {
		return "", errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespace)
	if !s.store.NamespaceExists(ns) {
		return "", errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	if version < 1 {
		return "", errors.New("[INVALID REQUEST]: Version must be greater than zero")
	}

	fullKey := ns + "/" + key

	versionKey := storage.VersionedSecretKey(fullKey, version)
	record, err := s.store.Get(versionKey)
	if err != nil {
		return "", errors.New("[SECRET NOT FOUND]: Requested version does not exist")
	}

	k := s.sealer.GetKey()
	plaintext, err := crypto.Decrypt(k, record.Ciphertext, record.Nonce)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
func (s *SecretService) ListVersions(namespaceName string, key string) ([]SecretVersion, error) {

	if s.sealer.IsSealed() {
		return nil, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	ns := normalizeNamespace(namespaceName)
	if !s.store.NamespaceExists(ns) {
		return nil, errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	fullKey := ns + "/" + key
	current, err := s.store.Get(fullKey)
	if err != nil {
		return nil, errors.New("[OPERATION DENIED]: Requested version does not exist")
	}

	var versions []SecretVersion
	for version := 1; version <= current.Version; version++ {
		versionKey := storage.VersionedSecretKey(fullKey, version)
		record, err := s.store.Get(versionKey)
		if err != nil {
			continue
		}

		versions = append(versions, SecretVersion{
			Version:   record.Version,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
			RotatedAt: record.RotatedAt,
		})

		log.Printf(
			"VERSION DEBUG -> version=%d created=%q updated=%q rotated=%q",
			record.Version,
			record.CreatedAt,
			record.UpdatedAt,
			record.RotatedAt,
		)
	}

	return versions, nil

}
