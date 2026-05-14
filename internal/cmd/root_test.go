package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd_Help(t *testing.T) {
	cmd := newRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	require.NoError(t, err)
	out := buf.String()
	require.True(t, strings.Contains(out, "envseal") || strings.Contains(out, "Usage"))
}

func TestRootCmd_UnknownSubcommand(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"nonexistent-subcommand"})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestRootCmd_HasExpectedSubcommands(t *testing.T) {
	cmd := newRootCmd()
	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	expected := []string{"seal", "open", "diff", "rotate", "keys", "version", "init"}
	for _, name := range expected {
		require.True(t, names[name], "expected subcommand %q to be registered", name)
	}
}

func TestKeystoreDir_Default(t *testing.T) {
	dir := keystoreDir("")
	require.NotEmpty(t, dir)
}

func TestStoreDir_Default(t *testing.T) {
	dir := storeDir("")
	require.NotEmpty(t, dir)
}
