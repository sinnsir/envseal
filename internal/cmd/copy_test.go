package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyCmd_NoArgs(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"copy"})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestCopyCmd_SameEnv(t *testing.T) {
	dir := t.TempDir()
	cmd := newRootCmd()
	out := new(bytes.Buffer)
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{
		"--keystore-dir", filepath.Join(dir, "keys"),
		"--store-dir", filepath.Join(dir, "sealed"),
		"copy", "production", "production",
	})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "must differ")
}

func TestCopyCmd_MissingSourceKey(t *testing.T) {
	dir := t.TempDir()
	cmd := newRootCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{
		"--keystore-dir", filepath.Join(dir, "keys"),
		"--store-dir", filepath.Join(dir, "sealed"),
		"copy", "production", "staging",
	})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestCopyCmd_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	keysDir := filepath.Join(dir, "keys")
	sealedDir := filepath.Join(dir, "sealed")
	require.NoError(t, os.MkdirAll(keysDir, 0700))
	require.NoError(t, os.MkdirAll(sealedDir, 0700))

	// Init source environment
	initCmd := newRootCmd()
	initCmd.SetOut(new(bytes.Buffer))
	initCmd.SetErr(new(bytes.Buffer))
	initCmd.SetArgs([]string{
		"--keystore-dir", keysDir,
		"--store-dir", sealedDir,
		"init", "production",
	})
	require.NoError(t, initCmd.Execute())

	// Write a .env file and seal it
	envFile := filepath.Join(dir, ".env.production")
	require.NoError(t, os.WriteFile(envFile, []byte("KEY=value\nSECRET=abc123\n"), 0600))

	sealCmd := newRootCmd()
	sealCmd.SetOut(new(bytes.Buffer))
	sealCmd.SetErr(new(bytes.Buffer))
	sealCmd.SetArgs([]string{
		"--keystore-dir", keysDir,
		"--store-dir", sealedDir,
		"seal", "production", envFile,
	})
	require.NoError(t, sealCmd.Execute())

	// Copy production → staging
	out := new(bytes.Buffer)
	copyCmd := newRootCmd()
	copyCmd.SetOut(out)
	copyCmd.SetErr(new(bytes.Buffer))
	copyCmd.SetArgs([]string{
		"--keystore-dir", keysDir,
		"--store-dir", sealedDir,
		"copy", "production", "staging",
	})
	require.NoError(t, copyCmd.Execute())
	require.Contains(t, out.String(), "production")
	require.Contains(t, out.String(), "staging")

	// Open staging and verify contents
	openOut := new(bytes.Buffer)
	openCmd := newRootCmd()
	openCmd.SetOut(openOut)
	openCmd.SetErr(new(bytes.Buffer))
	openCmd.SetArgs([]string{
		"--keystore-dir", keysDir,
		"--store-dir", sealedDir,
		"export", "staging", "--format", "dotenv",
	})
	require.NoError(t, openCmd.Execute())
	require.Contains(t, openOut.String(), "KEY=value")
	require.Contains(t, openOut.String(), "SECRET=abc123")
}
