package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

// ErrNotFound is returned when a key for the given environment does not exist.
var ErrNotFound = errors.New("key not found")

// KeyStore manages age identity files per environment.
type KeyStore struct {
	dir string
}

// New creates a KeyStore rooted at dir, creating it if necessary.
func New(dir string) (*KeyStore, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("create keystore dir: %w", err)
	}
	return &KeyStore{dir: dir}, nil
}

// Generate creates a new age X25519 identity for env, persisting it to disk.
// Any existing key is overwritten.
func (ks *KeyStore) Generate(env string) (*age.X25519Identity, error) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("generate identity: %w", err)
	}
	path := ks.path(env)
	if err := os.WriteFile(path, []byte(id.String()+"\n"), 0600); err != nil {
		return nil, fmt.Errorf("write key file: %w", err)
	}
	return id, nil
}

// Load reads the age identity for env from disk.
func (ks *KeyStore) Load(env string) (*age.X25519Identity, error) {
	path := ks.path(env)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, env)
	}
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}
	ids, err := age.ParseIdentities(bytesReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse identity: %w", err)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no identities found in key file for %q", env)
	}
	id, ok := ids[0].(*age.X25519Identity)
	if !ok {
		return nil, fmt.Errorf("unexpected identity type for %q", env)
	}
	return id, nil
}

// Exists reports whether a key exists for env.
func (ks *KeyStore) Exists(env string) bool {
	_, err := os.Stat(ks.path(env))
	return err == nil
}

// Delete removes the key for env. Returns ErrNotFound if it does not exist.
func (ks *KeyStore) Delete(env string) error {
	path := ks.path(env)
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %s", ErrNotFound, env)
	}
	return err
}

// List returns all environment names that have keys stored.
func (ks *KeyStore) List() ([]string, error) {
	entries, err := os.ReadDir(ks.dir)
	if err != nil {
		return nil, fmt.Errorf("read keystore dir: %w", err)
	}
	var envs []string
	for _, e := range entries {
		if !e.IsDir() {
			envs = append(envs, e.Name())
		}
	}
	return envs, nil
}

func (ks *KeyStore) path(env string) string {
	return filepath.Join(ks.dir, env)
}
