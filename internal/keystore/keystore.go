package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

const (
	KeyDir      = ".envseal"
	KeyFileExt  = ".key"
	PubFileExt  = ".pub"
)

// ErrKeyNotFound is returned when no key exists for the given environment.
var ErrKeyNotFound = errors.New("key not found for environment")

// KeyPair holds an age identity (private key) and its corresponding recipient.
type KeyPair struct {
	Environment string
	Identity    age.Identity
	Recipient   age.Recipient
}

// Store manages per-environment age key pairs on disk.
type Store struct {
	BaseDir string
}

// New creates a Store rooted at baseDir (usually the project root).
func New(baseDir string) *Store {
	return &Store{BaseDir: baseDir}
}

// Generate creates a new age X25519 key pair for the given environment and
// persists it under <BaseDir>/.envseal/<env>.key and <env>.pub.
func (s *Store) Generate(env string) (*KeyPair, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, fmt.Errorf("generating key: %w", err)
	}

	if err := os.MkdirAll(s.keyDir(), 0700); err != nil {
		return nil, fmt.Errorf("creating key directory: %w", err)
	}

	privPath := s.keyPath(env)
	if err := os.WriteFile(privPath, []byte(identity.String()+"\n"), 0600); err != nil {
		return nil, fmt.Errorf("writing private key: %w", err)
	}

	pubPath := s.pubPath(env)
	if err := os.WriteFile(pubPath, []byte(identity.Recipient().String()+"\n"), 0644); err != nil {
		return nil, fmt.Errorf("writing public key: %w", err)
	}

	return &KeyPair{
		Environment: env,
		Identity:    identity,
		Recipient:   identity.Recipient(),
	}, nil
}

// Load reads the age identity for the given environment from disk.
func (s *Store) Load(env string) (*KeyPair, error) {
	data, err := os.ReadFile(s.keyPath(env))
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("reading private key: %w", err)
	}

	identities, err := age.ParseIdentities(bytesReader(data))
	if err != nil || len(identities) == 0 {
		return nil, fmt.Errorf("parsing identity for %q: %w", env, err)
	}

	id := identities[0].(*age.X25519Identity)
	return &KeyPair{
		Environment: env,
		Identity:    id,
		Recipient:   id.Recipient(),
	}, nil
}

func (s *Store) keyDir() string  { return filepath.Join(s.BaseDir, KeyDir) }
func (s *Store) keyPath(e string) string { return filepath.Join(s.keyDir(), e+KeyFileExt) }
func (s *Store) pubPath(e string) string { return filepath.Join(s.keyDir(), e+PubFileExt) }
