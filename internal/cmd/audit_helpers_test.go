package cmd

import (
	"testing"

	"github.com/nicholasgasior/envseal/internal/envelope"
	"github.com/nicholasgasior/envseal/internal/keystore"
	"github.com/nicholasgasior/envseal/internal/store"
)

// sealEnvForTest is a shared helper that generates a key, seals vars, and
// writes the sealed blob into the store — used by multiple audit/status tests.
func sealEnvForTest(t *testing.T, ks *keystore.Keystore, st *store.Store, env string, vars map[string]string) error {
	t.Helper()

	identity, err := ks.Generate(env)
	if err != nil {
		return err
	}

	recipient, err := identity.Recipient()
	if err != nil {
		return err
	}

	sealed, err := envelope.Seal(env, vars, recipient)
	if err != nil {
		return err
	}

	return st.Write(env, sealed)
}

// newTestKeystoreAndStore creates isolated keystore and store instances backed
// by temporary directories, suitable for use in a single test.
func newTestKeystoreAndStore(t *testing.T) (*keystore.Keystore, *store.Store) {
	t.Helper()
	ksDir := t.TempDir()
	stDir := t.TempDir()
	ks, err := keystore.New(ksDir)
	if err != nil {
		t.Fatalf("keystore.New: %v", err)
	}
	st, err := store.New(stDir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return ks, st
}
