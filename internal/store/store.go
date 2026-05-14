// Package store manages sealed envelope files on disk.
package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const sealedExt = ".sealed"

// Store manages sealed .env files in a directory.
type Store struct {
	dir string
}

// New returns a Store rooted at dir, creating it if necessary.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Default returns the default store path relative to dir.
func Default(base string) string {
	return filepath.Join(base, ".envseal")
}

func (s *Store) path(env string) string {
	return filepath.Join(s.dir, env+sealedExt)
}

// Write persists the sealed data for env.
func (s *Store) Write(env string, data []byte) error {
	if err := os.WriteFile(s.path(env), data, 0o600); err != nil {
		return fmt.Errorf("write sealed env: %w", err)
	}
	return nil
}

// Read loads the sealed data for env.
func (s *Store) Read(env string) ([]byte, error) {
	data, err := os.ReadFile(s.path(env))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("sealed env %q not found", env)
		}
		return nil, fmt.Errorf("read sealed env: %w", err)
	}
	return data, nil
}

// Delete removes the sealed file for env.
func (s *Store) Delete(env string) error {
	if err := os.Remove(s.path(env)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("sealed env %q not found", env)
		}
		return fmt.Errorf("delete sealed env: %w", err)
	}
	return nil
}

// Exists reports whether a sealed file for env exists.
func (s *Store) Exists(env string) bool {
	_, err := os.Stat(s.path(env))
	return err == nil
}

// List returns all environment names that have sealed files.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("list store: %w", err)
	}
	var envs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), sealedExt) {
			envs = append(envs, strings.TrimSuffix(e.Name(), sealedExt))
		}
	}
	return envs, nil
}
