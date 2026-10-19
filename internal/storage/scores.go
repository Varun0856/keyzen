package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Best struct {
	WPM      float64   `json:"wpm"`
	Accuracy float64   `json:"accuracy"`
	Date     time.Time `json:"date"`
}

type Scores map[string]Best

type Store struct {
	path string
}

func New(path string) *Store { return &Store{path: path} }

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "keyzen", "scores.json"), nil
}

func key(seconds int) string { return strconv.Itoa(seconds) + "s" }

func (s *Store) Load() (Scores, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return Scores{}, nil
	}
	if err != nil {
		return nil, err
	}

	var sc Scores
	if err := json.Unmarshal(data, &sc); err != nil {
		return nil, fmt.Errorf("parse %s: %w", s.path, err)
	}
	if sc == nil {
		sc = Scores{}
	}
	return sc, nil
}

func (s *Store) Save(sc Scores) error {
	data, err := json.MarshalIndent(sc, "", " ")
	if err != nil {
		return err
	}
	return writeAtomic(s.path, data)
}

func (s *Store) Record(key string, wpm, accuracy float64) (Best, bool, error) {
	sc, err := s.Load()
	if err != nil {
		return Best{}, false, err
	}

	prev, ok := sc[key]
	if ok && wpm <= prev.WPM {
		return prev, false, nil
	}

	best := Best{WPM: wpm, Accuracy: accuracy, Date: time.Now().UTC()}
	sc[key] = best
	if err := s.Save(sc); err != nil {
		return Best{}, false, err
	}
	return best, true, nil
}
