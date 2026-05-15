package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRotateCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"rotate"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	err := root.Execute()
	require.Error(t, err)
}

func TestRotateCmd_MissingEnv(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", filepath.Join(tmp, "keys"))
	t.Setenv("ENVSEAL_STORE", filepath.Join(tmp, "store"))

	root := newRootCmd()
	root.SetArgs([]string{"rotate", "production"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	err := root.Execute()
	require.Error(t, err)
}

func TestRotateCmd_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", filepath.Join(tmp, "keys"))
	t.Setenv("ENVSEAL_STORE", filepath.Join(tmp, "store"))
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "store"), 0700))

	plainFile := filepath.Join(tmp, ".env.staging")
	require.NoError(t, os.WriteFile(plainFile, []byte("KEY=value\nFOO=bar\n"), 0600))

	root := newRootCmd()
	root.SetArgs([]string{"seal", "--env", "staging", plainFile})
	require.NoError(t, root.Execute())

	root2 := newRootCmd()
	root2.SetArgs([]string{"rotate", "staging"})
	var out bytes.Buffer
	root2.SetOut(&out)
	require.NoError(t, root2.Execute())
	require.Contains(t, out.String(), "staging")

	outFile := filepath.Join(tmp, "decrypted.env")
	root3 := newRootCmd()
	root3.SetArgs([]string{"open", "--env", "staging", "--output", outFile})
	require.NoError(t, root3.Execute())

	got, err := os.ReadFile(outFile)
	require.NoError(t, err)
	require.Contains(t, string(got), "KEY=value")
}

// TestRotateCmd_Idempotent verifies that rotating the same environment twice
// succeeds and still produces a valid, decryptable store.
func TestRotateCmd_Idempotent(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", filepath.Join(tmp, "keys"))
	t.Setenv("ENVSEAL_STORE", filepath.Join(tmp, "store"))
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, "store"), 0700))

	plainFile := filepath.Join(tmp, ".env.prod")
	require.NoError(t, os.WriteFile(plainFile, []byte("SECRET=abc123\n"), 0600))

	// Initial seal.
	root := newRootCmd()
	root.SetArgs([]string{"seal", "--env", "prod", plainFile})
	require.NoError(t, root.Execute())

	// First rotation.
	root2 := newRootCmd()
	root2.SetArgs([]string{"rotate", "prod"})
	require.NoError(t, root2.Execute())

	// Second rotation.
	root3 := newRootCmd()
	root3.SetArgs([]string{"rotate", "prod"})
	require.NoError(t, root3.Execute())

	// Verify the data is still readable after two rotations.
	outFile := filepath.Join(tmp, "decrypted.env")
	root4 := newRootCmd()
	root4.SetArgs([]string{"open", "--env", "prod", "--output", outFile})
	require.NoError(t, root4.Execute())

	got, err := os.ReadFile(outFile)
	require.NoError(t, err)
	require.Contains(t, string(got), "SECRET=abc123")
}
