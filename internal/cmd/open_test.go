package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenCmd_NoArgs(t *testing.T) {
	cmd := newOpenCmd(t.TempDir(), t.TempDir())
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestOpenCmd_MissingSealedEnv(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	// Init key so keystore exists
	initCmd := newInitCmd(ksDir, stDir)
	initCmd.SetArgs([]string{"production"})
	err := initCmd.Execute()
	require.NoError(t, err)

	// Try to open env that was never sealed
	openCmd := newOpenCmd(ksDir, stDir)
	out := filepath.Join(t.TempDir(), ".env.out")
	openCmd.SetArgs([]string{"production", "--output", out})
	err = openCmd.Execute()
	require.Error(t, err)
}

func TestOpenCmd_MissingKey(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	// Write a fake sealed file without a key
	err := os.MkdirAll(stDir, 0700)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(stDir, "staging.age"), []byte("not-real-ciphertext"), 0600)
	require.NoError(t, err)

	openCmd := newOpenCmd(ksDir, stDir)
	out := filepath.Join(t.TempDir(), ".env.out")
	openCmd.SetArgs([]string{"staging", "--output", out})
	err = openCmd.Execute()
	require.Error(t, err)
}
