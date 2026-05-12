package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envseal/internal/keystore"
)

func TestGenerateAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	store := keystore.New(tmpDir)

	kp, err := store.Generate("production")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if kp.Environment != "production" {
		t.Errorf("expected environment %q, got %q", "production", kp.Environment)
	}
	if kp.Identity == nil || kp.Recipient == nil {
		t.Fatal("expected non-nil Identity and Recipient")
	}

	// Key files should exist on disk.
	for _, name := range []string{"production.key", "production.pub"} {
		path := filepath.Join(tmpDir, keystore.KeyDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", path, err)
		}
	}

	// Private key file should not be world-readable.
	info, _ := os.Stat(filepath.Join(tmpDir, keystore.KeyDir, "production.key"))
	if info.Mode().Perm() != 0600 {
		t.Errorf("private key perm: got %v, want 0600", info.Mode().Perm())
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	store := keystore.New(tmpDir)

	orig, err := store.Generate("staging")
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	loaded, err := store.Load("staging")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if orig.Recipient.String() != loaded.Recipient.String() {
		t.Errorf("recipient mismatch after round-trip")
	}
}

func TestLoad_NotFound(t *testing.T) {
	store := keystore.New(t.TempDir())
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != keystore.ErrKeyNotFound {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}
