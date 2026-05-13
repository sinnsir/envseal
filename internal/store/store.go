package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultDir      = ".envseal"
	SealedExtension = ".sealed"
)

// ErrNotFound is returned when a sealed file does not exist for the given environment.
var ErrNotFound = errors.New("sealed env file not found")

// Store manages sealed .env files on disk.
type Store struct {
	baseDir string
}

// New creates a Store rooted at baseDir.
func New(baseDir string) *Store {
	return &Store{baseDir: baseDir}
}

// Default returns a Store using the default directory relative to dir.
func Default(dir string) *Store {
	return New(filepath.Join(dir, DefaultDir))
}

// Path returns the file path for the given environment's sealed file.
func (s *Store) Path(env string) string {
	return filepath.Join(s.baseDir, env+SealedExtension)
}

// Write persists data as the sealed file for env, creating directories as needed.
func (s *Store) Write(env string, data []byte) error {
	if err := os.MkdirAll(s.baseDir, 0o700); err != nil {
		return fmt.Errorf("store: create dir: %w", err)
	}
	path := s.Path(env)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("store: write %s: %w", path, err)
	}
	return nil
}

// Read returns the sealed bytes for env.
func (s *Store) Read(env string) ([]byte, error) {
	path := s.Path(env)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("store: %w for environment %q", ErrNotFound, env)
	}
	if err != nil {
		return fmt.Errorf("store: read %s: %w", path, err)
	}
	return data, nil
}

// List returns all environment names that have sealed files.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.baseDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("store: list: %w", err)
	}
	var envs []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) == SealedExtension {
			envs = append(envs, name[:len(name)-len(SealedExtension)])
		}
	}
	return envs, nil
}

// Remove deletes the sealed file for env.
func (s *Store) Remove(env string) error {
	path := s.Path(env)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("store: %w for environment %q", ErrNotFound, env)
		}
		return fmt.Errorf("store: remove %s: %w", path, err)
	}
	return nil
}
