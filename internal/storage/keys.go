package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Varun0856/keyzen/internal/typing"
)

type KeyTotals map[string]typing.KeyStat

type KeyStore struct{ path string }

func NewKeyStore(path string) *KeyStore { return &KeyStore{path: path} }

func DefaultKeysPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "keyzen", "keys.json"), nil
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *KeyStore) Load() (KeyTotals, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return KeyTotals{}, nil
	}
	if err != nil {
		return nil, err
	}

	var t KeyTotals
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parse %s: %w", s.path, err)
	}
	if t == nil {
		t = KeyTotals{}
	}
	return t, nil
}

func (t KeyTotals) Merge(session map[rune]typing.KeyStat) {
	for r, st := range session {
		key := string(r)
		t[key] = t[key].Add(st)
	}
}

func (t KeyTotals) Runes() map[rune]typing.KeyStat {
	out := make(map[rune]typing.KeyStat, len(t))
	for k, st := range t {
		if r := []rune(k); len(r) == 1 {
			out[r[0]] = st
		}
	}
	return out
}

func (s *KeyStore) Record(session map[rune]typing.KeyStat) error {
	t, err := s.Load()
	if err != nil {
		return err
	}
	t.Merge(session)

	data, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}
	return writeAtomic(s.path, data)
}
