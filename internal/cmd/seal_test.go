package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSealCmd_NoArgs(t *testing.T) {
	cmd := newSealCmd(t.TempDir(), t.TempDir())
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestSealCmd_MissingEnvFile(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()
	cmd := newSealCmd(ksDir, stDir)
	cmd.SetArgs([]string{"production", "--file", "/nonexistent/.env"})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestSealCmd_RoundTrip(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	// Write a sample .env file
	envFile := filepath.Join(t.TempDir(), ".env")
	err := os.WriteFile(envFile, []byte("FOO=bar\nBAZ=qux\n"), 0600)
	require.NoError(t, err)

	// Init key first
	initCmd := newInitCmd(ksDir, stDir)
	initCmd.SetArgs([]string{"staging"})
	err = initCmd.Execute()
	require.NoError(t, err)

	// Seal
	sealCmd := newSealCmd(ksDir, stDir)
	sealCmd.SetArgs([]string{"staging", "--file", envFile})
	err = sealCmd.Execute()
	require.NoError(t, err)

	// Verify sealed file exists
	openCmd := newOpenCmd(ksDir, stDir)
	out := filepath.Join(t.TempDir(), ".env.out")
	openCmd.SetArgs([]string{"staging", "--output", out})
	err = openCmd.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(out)
	require.NoError(t, err)
	require.Contains(t, string(data), "FOO=bar")
	require.Contains(t, string(data), "BAZ=qux")
}
