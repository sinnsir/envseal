package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd_Help(t *testing.T) {
	out := new(bytes.Buffer)
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"--help"})
	_ = root.Execute()
	if !strings.Contains(out.String(), "envseal") {
		t.Errorf("expected help to mention envseal, got: %s", out.String())
	}
}

func TestRootCmd_UnknownSubcommand(t *testing.T) {
	root := newRootCmd()
	root.SetOut(new(bytes.Buffer))
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"notacommand"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unknown subcommand")
	}
}

func TestRootCmd_HasExpectedSubcommands(t *testing.T) {
	root := newRootCmd()
	names := make(map[string]bool)
	for _, c := range root.Commands() {
		names[c.Name()] = true
	}
	expected := []string{"init", "seal", "open", "diff", "rotate", "keys", "version", "edit", "export"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestKeystoreDir_Default(t *testing.T) {
	t.Setenv("ENVSEAL_KEYSTORE", "")
	dir := keystoreDir()
	if dir == "" {
		t.Error("expected non-empty keystore dir")
	}
	if strings.Contains(dir, "ENVSEAL") {
		t.Errorf("expected resolved path, got: %s", dir)
	}
}

func TestStoreDir_Default(t *testing.T) {
	t.Setenv("ENVSEAL_STORE", "")
	dir := storeDir()
	if dir == "" {
		t.Error("expected non-empty store dir")
	}
}
