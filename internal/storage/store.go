package storage

type Store interface {
	Put(key string, record SecretRecord) error
	Get(key string) (SecretRecord, error)
}

type SecretRecord struct {
	Ciphertext []byte
	Nonce      []byte
}

var _ Store = (*BoltStore)(nil)
