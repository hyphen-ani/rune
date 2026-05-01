package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Token string `json:"token"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".rune")

	// Create ~/.rune if it doesnt exists
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), nil

}

func SaveToken(token string) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	cfg := Config{
		Token: token,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)

}

func LoadToken() (string, error) {
	path, err := getConfigPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return "", err
	}

	return cfg.Token, nil

}
