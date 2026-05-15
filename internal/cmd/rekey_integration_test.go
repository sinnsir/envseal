package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

// TestRekey_OldKeyCannotDecrypt verifies that after rekeying, the old key
// identity is replaced and the old private key material is no longer stored.
func TestRekey_OldKeyCannotDecrypt(t *testing.T) {
	dir := t.TempDir()
	ksDir := filepath.Join(dir, "keys")
	stDir := filepath.Join(dir, "store")
	os.MkdirAll(ksDir, 0700) //nolint:errcheck
	os.MkdirAll(stDir, 0700) //nolint:errcheck

	// Init + seal.
	run := func(args ...string) {
		t.Helper()
		cmd := newRootCmd()
		cmd.SetArgs(append([]string{"--keystore-dir", ksDir, "--store-dir", stDir}, args...))
		if err := cmd.Execute(); err != nil {
			t.Fatalf("cmd %v: %v", args, err)
		}
	}

	run("init", "prod")

	tmpEnv := filepath.Join(dir, ".env")
	os.WriteFile(tmpEnv, []byte("SECRET=hunter2\n"), 0600) //nolint:errcheck
	run("seal", "prod", tmpEnv)

	// Capture the old identity before rekeying.
	ks, _ := keystore.New(ksDir)
	oldID, err := ks.Load("prod")
	if err != nil {
		t.Fatalf("load old key: %v", err)
	}

	// Rekey.
	run("rekey", "prod")

	// The sealed file should now be decryptable only with the new key.
	newID, err := ks.Load("prod")
	if err != nil {
		t.Fatalf("load new key: %v", err)
	}

	st, _ := store.New(stDir)
	sealed, err := st.Read("prod")
	if err != nil {
		t.Fatalf("read sealed: %v", err)
	}

	// New key must succeed.
	if _, err := envelope.Open(sealed, newID); err != nil {
		t.Errorf("new key should decrypt: %v", err)
	}

	// Old key must fail.
	if _, err := envelope.Open(sealed, oldID); err == nil {
		t.Error("old key should NOT decrypt after rekey")
	}
}
