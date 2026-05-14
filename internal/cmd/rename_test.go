package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenameCmd_NoArgs(t *testing.T) {
	cmd := newRenameCmd(func() string { return t.TempDir() }, func() string { return t.TempDir() })
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestRenameCmd_SameEnv(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()
	cmd := newRenameCmd(func() string { return ksDir }, func() string { return stDir })
	cmd.SetArgs([]string{"prod", "prod"})
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "must differ")
}

func TestRenameCmd_MissingSourceKey(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()
	cmd := newRenameCmd(func() string { return ksDir }, func() string { return stDir })
	cmd.SetArgs([]string{"nonexistent", "newenv"})
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no key found")
}

func TestRenameCmd_RoundTrip(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	// Init source environment
	initCmd := newInitCmd(func() string { return ksDir })
	initCmd.SetArgs([]string{"staging"})
	initCmd.SetOut(&bytes.Buffer{})
	initCmd.SetErr(&bytes.Buffer{})
	require.NoError(t, initCmd.Execute())

	// Seal a file
	sealCmd := newSealCmd(func() string { return ksDir }, func() string { return stDir })
	envFile := writeEnvFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	sealCmd.SetArgs([]string{"staging", envFile})
	sealCmd.SetOut(&bytes.Buffer{})
	sealCmd.SetErr(&bytes.Buffer{})
	require.NoError(t, sealCmd.Execute())

	// Rename staging -> production
	renameCmd := newRenameCmd(func() string { return ksDir }, func() string { return stDir })
	renameCmd.SetArgs([]string{"staging", "production"})
	out := &bytes.Buffer{}
	renameCmd.SetOut(out)
	renameCmd.SetErr(&bytes.Buffer{})
	require.NoError(t, renameCmd.Execute())
	require.Contains(t, out.String(), "renamed")

	// Old env should be gone
	ksOld, err := newKeystore(ksDir)
	require.NoError(t, err)
	require.False(t, ksOld.Exists("staging"))

	// New env should be openable
	openCmd := newOpenCmd(func() string { return ksDir }, func() string { return stDir })
	openOut := &bytes.Buffer{}
	openCmd.SetArgs([]string{"production"})
	openCmd.SetOut(openOut)
	openCmd.SetErr(&bytes.Buffer{})
	require.NoError(t, openCmd.Execute())
	require.Contains(t, openOut.String(), "DB_HOST")
}
