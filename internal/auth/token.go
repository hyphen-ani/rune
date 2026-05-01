package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateToken() (string, string) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
	}

	token := hex.EncodeToString(b)

	hash := sha256.Sum256([]byte(token))
	hashStr := hex.EncodeToString(hash[:])

	return token, hashStr
}
