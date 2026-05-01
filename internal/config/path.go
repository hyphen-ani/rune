package config

import (
	"os"
	"path/filepath"
)

func GetRuneDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".rune")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	return dir, nil
}

func GetDBPath() (string, error) {
	dir, err := GetRuneDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "rune.db"), nil
}
