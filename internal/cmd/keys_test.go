package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKeysListCmd_NoKeys(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	root := newRootCmd()
	root.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "keys", "list"})

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "no keys") {
		t.Errorf("expected 'no keys' message, got: %q", out)
	}
}

func TestKeysListCmd_WithKey(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	// Create a key by running init for an environment
	root := newRootCmd()
	root.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "init", "--env", "staging"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Now list keys
	root2 := newRootCmd()
	root2.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "keys", "list"})
	var buf2 bytes.Buffer
	root2.SetOut(&buf2)
	root2.SetErr(&buf2)
	if err := root2.Execute(); err != nil {
		t.Fatalf("keys list failed: %v", err)
	}

	out := buf2.String()
	if !strings.Contains(out, "staging") {
		t.Errorf("expected 'staging' in output, got: %q", out)
	}
}

func TestKeysDeleteCmd_Existing(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	// Init an environment to create a key
	root := newRootCmd()
	root.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "init", "--env", "dev"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	// Delete the key
	root2 := newRootCmd()
	root2.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "keys", "delete", "--env", "dev"})
	var buf2 bytes.Buffer
	root2.SetOut(&buf2)
	root2.SetErr(&buf2)
	if err := root2.Execute(); err != nil {
		t.Fatalf("keys delete failed: %v", err)
	}

	// Verify key file is gone
	keyFile := filepath.Join(ksDir, "dev.key")
	if _, err := os.Stat(keyFile); !os.IsNotExist(err) {
		t.Errorf("expected key file to be deleted")
	}
}
