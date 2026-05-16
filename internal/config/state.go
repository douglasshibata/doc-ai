package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type State struct {
	Files map[string]string `json:"files"` // map[filePath]hash
}

func LoadState(root string) (*State, error) {
	statePath := filepath.Join(root, ".docai", "state.json")
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return &State{Files: make(map[string]string)}, nil
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *State) Save(root string) error {
	dir := filepath.Join(root, ".docai")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	statePath := filepath.Join(dir, "state.json")
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0644)
}

func CalculateHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
