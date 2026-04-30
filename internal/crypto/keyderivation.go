package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

func DeriveKey(passphrase string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		1,
		64*1024,
		4,
		32,
	)
}
