package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecret() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil

}
