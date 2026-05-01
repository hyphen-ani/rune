package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type TokenRecord struct {
	ID        string `json:"id"`
	Hash      string `json:"hash"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Revoked   bool   `json:"revoked"`
}

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
