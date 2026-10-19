package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Theme string `json:"theme"`
}

type ConfigStore struct{ path string }

func NewConfigStore(path string) *ConfigStore { return &ConfigStore{path: path} }

func DefaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "keyzen", "config.json"), nil
}

func (s *ConfigStore) Load() (Config, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}

	if err != nil {
		return Config{}, err
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", s.path, err)
	}
	return c, nil
}

func (s *ConfigStore) Save(c Config) error {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	return writeAtomic(s.path, data)
}
