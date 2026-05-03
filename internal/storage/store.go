package storage

import "rune/internal/auth"

type Store interface {
	Put(key string, record SecretRecord) error
	Get(key string) (SecretRecord, error)
	ListKeys() ([]string, error)

	GetSalt() ([]byte, error)
	SaveSalt(salt []byte) error

	SaveToken(record auth.TokenRecord) error
	GetTokenByHash(hash string) (*auth.TokenRecord, error)
	ListTokens() ([]auth.TokenRecord, error)
	RevokeToken(id string) error
	IsValidToken(hash string) bool

	CreateNamespace(namespace string) error
	NamespaceExists(namespace string) bool
	ListNamespaces() ([]string, error)
}

const VerifyKey = "__rune_verify"

type SecretRecord struct {
	Ciphertext []byte
	Nonce      []byte
}

var _ Store = (*BoltStore)(nil)
