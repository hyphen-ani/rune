package storage

import (
	"rune/internal/auth"
)

type Store interface {
	Put(key string, record SecretRecord) error
	Get(key string) (SecretRecord, error)
	Delete(key string) error
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
	DeleteNamespace(namespace string) error

	GetStats() (VaultStats, error)
}

type NamespaceCount struct {
	Namespace string `json:"namespace"`
	Count     int    `json:"count"`
}

type MonthCount struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

type VaultStats struct {
	TotalSecrets       int              `json:"total_secrets"`
	TotalNamespaces    int              `json:"total_namespaces"`
	RotatedSecrets     int              `json:"rotated_secrets"`
	TotalTokens        int              `json:"total_tokens"`
	ActiveTokens       int              `json:"active_tokens"`
	RevokedTokens      int              `json:"revoked_tokens"`
	SecretsByNamespace []NamespaceCount `json:"secrets_by_namespace"`
	TokensByMonth      []MonthCount     `json:"tokens_by_month"`
}

const VerifyKey = "__rune_verify"

type SecretRecord struct {
	Ciphertext []byte
	Nonce      []byte

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	RotatedAt string `json:"rotated_at"`
	Version   int    `json:"version"`
}

var _ Store = (*BoltStore)(nil)
