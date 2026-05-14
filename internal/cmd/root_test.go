package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd_Help(t *testing.T) {
	root := newRootCmd()
	out := &bytes.Buffer{}
	root.cmd.SetOut(out)
	root.cmd.SetErr(&bytes.Buffer{})
	root.cmd.SetArgs([]string{"--help"})
	err := root.cmd.Execute()
	require.NoError(t, err)
	require.Contains(t, out.String(), "envseal")
}

func TestRootCmd_UnknownSubcommand(t *testing.T) {
	root := newRootCmd()
	root.cmd.SetOut(&bytes.Buffer{})
	root.cmd.SetErr(&bytes.Buffer{})
	root.cmd.SetArgs([]string{"doesnotexist"})
	err := root.cmd.Execute()
	require.Error(t, err)
}

func TestRootCmd_HasExpectedSubcommands(t *testing.T) {
	root := newRootCmd()
	names := make([]string, 0)
	for _, sub := range root.cmd.Commands() {
		names = append(names, sub.Name())
	}
	expected := []string{"init", "seal", "open", "edit", "diff", "rotate", "copy", "rename", "export", "list", "status", "keys", "version"}
	for _, e := range expected {
		require.True(t, containsString(names, e), "expected subcommand %q", e)
	}
}

func TestKeystoreDir_Default(t *testing.T) {
	dir := keystoreDir()
	require.True(t, strings.Contains(dir, "envseal") || strings.Contains(dir, ".config") || len(dir) > 0)
}

func TestStoreDir_Default(t *testing.T) {
	dir := storeDir()
	require.NotEmpty(t, dir)
}

func containsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
