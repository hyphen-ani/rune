package storage

type Store interface {
	Put(key string, record SecretRecord) error
	Get(key string) (SecretRecord, error)
	ListKeys() ([]string, error)
	GetSalt() ([]byte, error)
	SaveSalt(salt []byte) error
	SaveToken(hash string) error
	IsValidToken(hash string) bool
}

const VerifyKey = "__rune_verify"

type SecretRecord struct {
	Ciphertext []byte
	Nonce      []byte
}

var _ Store = (*BoltStore)(nil)
