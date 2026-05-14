package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envseal/internal/keystore"
	"github.com/nicholasgasior/envseal/internal/store"
)

func TestListCmd_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("ENVSEAL_STORE_DIR", filepath.Join(tmpDir, "store"))
	t.Setenv("ENVSEAL_KEYSTORE_DIR", filepath.Join(tmpDir, "keys"))

	root := newRootCmd()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"list"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "No sealed environments found.") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestListCmd_WithEnvironments(t *testing.T) {
	tmpDir := t.TempDir()
	storeDir := filepath.Join(tmpDir, "store")
	keysDir := filepath.Join(tmpDir, "keys")
	t.Setenv("ENVSEAL_STORE_DIR", storeDir)
	t.Setenv("ENVSEAL_KEYSTORE_DIR", keysDir)

	ks, err := keystore.New(keysDir)
	if err != nil {
		t.Fatalf("keystore: %v", err)
	}
	identity, err := ks.Generate("production")
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	_ = identity

	envFile := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envFile, []byte("FOO=bar\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	s, err := store.New(storeDir)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	_ = s

	root := newRootCmd()
	root.SetArgs([]string{"seal", "-e", "production", envFile})
	if err := root.Execute(); err != nil {
		t.Fatalf("seal: %v", err)
	}

	buf := &bytes.Buffer{}
	root2 := newRootCmd()
	root2.SetOut(buf)
	root2.SetArgs([]string{"list"})
	if err := root2.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}

	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected 'production' in output, got: %s", buf.String())
	}
}
