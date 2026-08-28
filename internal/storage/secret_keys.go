package storage

import (
	"fmt"
	"strings"
)

func VersionedSecretKey(key string, version int) string {
	return fmt.Sprintf("%s@v%d", key, version)
}

func IsVersionedSecretKey(key string) bool {
	index := strings.LastIndex(key, "@v")
	if index == -1 {
		return false
	}

	if index+2 >= len(key) {
		return false
	}

	return true
}
