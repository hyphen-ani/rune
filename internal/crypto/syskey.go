package crypto

import "crypto/sha256"

// DeriveSysKey derives a 32-byte AES-256 key from a plain system token string.
// The domain prefix ensures the derived key is distinct from any other SHA256 use.
// Attackers who only have the DB (not the token) cannot derive this key.
func DeriveSysKey(plainToken string) []byte {
	hash := sha256.Sum256([]byte("rune-system-key:" + plainToken))
	return hash[:]
}
