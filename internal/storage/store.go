package storage

type Store interface {
	Put(key string, record SecretRecord) error
	Get(key string) (SecretRecord, error)
	GetSalt() ([]byte, error)
	SaveSalt(salt []byte) error
}

const VerifyKey = "__rune_verify"

type SecretRecord struct {
	Ciphertext []byte
	Nonce      []byte
}

var _ Store = (*BoltStore)(nil)
