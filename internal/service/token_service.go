package service

import (
	"encoding/hex"
	"errors"
	"fmt"
	"rune/internal/auth"
	"rune/internal/crypto"
	"rune/internal/seal"
	"rune/internal/storage"
	"strings"
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
func (s *TokenService) Create(name string, namespaces []auth.NamespacePermission) (string, auth.TokenRecord, error) {

	if s.sealer.IsSealed() {
		return "", auth.TokenRecord{}, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	if len(namespaces) == 0 {
		return "", auth.TokenRecord{}, errors.New("[INVALID]: At least one namespace is required")
	}

	// Validate all namespaces and permissions.
	for i, np := range namespaces {
		ns := normalizeNamespace(np.Namespace)
		namespaces[i].Namespace = ns
		if !s.store.NamespaceExists(ns) {
			return "", auth.TokenRecord{}, fmt.Errorf("[OPERATION DENIED]: Namespace %q does not exist", ns)
		}
		for _, p := range np.Permissions {
			if !auth.ValidPermissions[p] {
				return "", auth.TokenRecord{}, fmt.Errorf("[INVALID]: Unknown permission %q", p)
			}
		}
		if len(np.Permissions) == 0 {
			return "", auth.TokenRecord{}, fmt.Errorf("[INVALID]: Namespace %q has no permissions", ns)
		}
	}

	token, hash := auth.GenerateToken()
	id := fmt.Sprintf("tkn_%d", time.Now().UnixNano())

	record := auth.TokenRecord{
		ID:         id,
		Hash:       hash,
		Name:       name,
		Namespaces: namespaces,
		CreatedAt:  time.Now().Format(time.RFC3339),
		Revoked:    false,
	}

	err := s.store.SaveToken(record)
	return token, record, err
}
func (s *TokenService) CreateSystemToken(name, namespace string) (string, auth.TokenRecord, error) {
	if s.sealer.IsSealed() {
		return "", auth.TokenRecord{}, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}

	namespace = normalizeNamespace(namespace)

	if !s.store.NamespaceExists(namespace) {
		return "", auth.TokenRecord{}, errors.New("[OPERATION DENIED]: Namespace does not exist")
	}

	plainToken, hash := auth.GenerateToken()
	sysKey := crypto.DeriveSysKey(plainToken)
	vaultKey := s.sealer.GetKey()

	encCiphertext, encNonce, err := crypto.Encrypt(vaultKey, sysKey)
	if err != nil {
		return "", auth.TokenRecord{}, err
	}
	encSysKey := hex.EncodeToString(encCiphertext) + ":" + hex.EncodeToString(encNonce)

	id := fmt.Sprintf("stk_%d", time.Now().UnixNano())
	record := auth.TokenRecord{
		ID:        id,
		Hash:      hash,
		Name:      name,
		Namespace: namespace, // kept for ListSystemTokensForNamespace lookup
		Namespaces: []auth.NamespacePermission{
			{Namespace: namespace, Permissions: []string{auth.PermRead}},
		},
		CreatedAt: time.Now().Format(time.RFC3339),
		Revoked:   false,
		ReadOnly:  true,
		System:    true,
		EncSysKey: encSysKey,
	}

	if err := s.store.SaveToken(record); err != nil {
		return "", auth.TokenRecord{}, err
	}

	// Encrypt all existing namespace secrets under the system key so this
	// token can read them even when the vault is sealed.
	keys, err := s.store.ListKeys()
	if err != nil {
		return plainToken, record, nil
	}

	nsPrefix := namespace + "/"
	for _, fullKey := range keys {
		if !strings.HasPrefix(fullKey, nsPrefix) || storage.IsVersionedSecretKey(fullKey) {
			continue
		}
		rec, err := s.store.Get(fullKey)
		if err != nil {
			continue
		}
		plaintext, err := crypto.Decrypt(vaultKey, rec.Ciphertext, rec.Nonce)
		if err != nil {
			continue
		}
		ct, nonce, err := crypto.Encrypt(sysKey, plaintext)
		if err != nil {
			continue
		}
		sysRecord := storage.SecretRecord{
			Ciphertext: ct,
			Nonce:      nonce,
			CreatedAt:  rec.CreatedAt,
			UpdatedAt:  rec.UpdatedAt,
		}
		justKey := strings.TrimPrefix(fullKey, nsPrefix)
		s.store.PutSystemSecret(id, namespace, justKey, sysRecord)
	}

	return plainToken, record, nil
}

func (s *TokenService) Revoke(id string) error {
	if s.sealer.IsSealed() {
		return errors.New("[OPERATION DENIED]: Vault is Sealed")
	}
	if id == "root" {
		return errors.New("[OPERATION DENIED]: Cannot revoke root token")
	}
	tokens, err := s.store.ListTokens()
	if err == nil {
		for _, t := range tokens {
			if t.ID == id && t.System {
				s.store.DeleteAllSystemSecrets(id)
				break
			}
		}
	}
	return s.store.RevokeToken(id)
}
func (s *TokenService) List() ([]auth.TokenRecord, error) {
	if s.sealer.IsSealed() {
		return nil, errors.New("[OPERATION DENIED]: Vault is Sealed")
	}
	return s.store.ListTokens()
}
