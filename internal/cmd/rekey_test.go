package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func TestRekeyCmd_NoArgs(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"rekey"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing argument")
	}
}

func TestRekeyCmd_MissingKey(t *testing.T) {
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"--keystore-dir", filepath.Join(dir, "keys"),
		"--store-dir", filepath.Join(dir, "store"),
		"rekey", "staging",
	})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when key does not exist")
	}
}

func TestRekeyCmd_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	ksDir := filepath.Join(dir, "keys")
	stDir := filepath.Join(dir, "store")
	if err := os.MkdirAll(ksDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(stDir, 0700); err != nil {
		t.Fatal(err)
	}

	// Bootstrap: init environment.
	init := newRootCmd()
	init.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "init", "staging"})
	if err := init.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Seal a file.
	tmpEnv := filepath.Join(dir, ".env")
	if err := os.WriteFile(tmpEnv, []byte("FOO=bar\nBAZ=qux\n"), 0600); err != nil {
		t.Fatal(err)
	}
	sealCmd := newRootCmd()
	sealCmd.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "seal", "staging", tmpEnv})
	if err := sealCmd.Execute(); err != nil {
		t.Fatalf("seal: %v", err)
	}

	// Rekey.
	rekeyCmd := newRootCmd()
	var out bytes.Buffer
	rekeyCmd.SetOut(&out)
	rekeyCmd.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "rekey", "staging"})
	if err := rekeyCmd.Execute(); err != nil {
		t.Fatalf("rekey: %v", err)
	}
	if !strings.Contains(out.String(), "rekeyed environment") {
		t.Errorf("unexpected output: %q", out.String())
	}

	// Verify the new key can decrypt the re-encrypted file.
	ks, _ := keystore.New(ksDir)
	st, _ := store.New(stDir)
	id, err := ks.Load("staging")
	if err != nil {
		t.Fatalf("load new key: %v", err)
	}
	sealed, err := st.Read("staging")
	if err != nil {
		t.Fatalf("read sealed: %v", err)
	}
	env_map, err := envelope.Open(sealed, id)
	if err != nil {
		t.Fatalf("open with new key: %v", err)
	}
	if env_map["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env_map["FOO"])
	}
}
