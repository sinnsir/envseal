package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatusCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"status"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestStatusCmd_MissingSealedEnv(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("ENVSEAL_STORE_DIR", filepath.Join(tmpDir, "store"))
	t.Setenv("ENVSEAL_KEYSTORE_DIR", filepath.Join(tmpDir, "keys"))

	root := newRootCmd()
	root.SetArgs([]string{"status", "staging"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing sealed env")
	}
}

func TestStatusCmd_NoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("ENVSEAL_STORE_DIR", filepath.Join(tmpDir, "store"))
	t.Setenv("ENVSEAL_KEYSTORE_DIR", filepath.Join(tmpDir, "keys"))

	envFile := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envFile, []byte("FOO=bar\nBAZ=qux\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	sealRoot := newRootCmd()
	sealRoot.SetArgs([]string{"init", "production"})
	if err := sealRoot.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	sealRoot2 := newRootCmd()
	sealRoot2.SetArgs([]string{"seal", "-e", "production", envFile})
	if err := sealRoot2.Execute(); err != nil {
		t.Fatalf("seal: %v", err)
	}

	buf := &bytes.Buffer{}
	statusRoot := newRootCmd()
	statusRoot.SetOut(buf)
	statusRoot.SetArgs([]string{"status", "-f", envFile, "production"})
	if err := statusRoot.Execute(); err != nil {
		t.Fatalf("status: %v", err)
	}

	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}
