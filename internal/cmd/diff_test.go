package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envseal/internal/cmd"
)

func TestDiffCmd_NoArgs(t *testing.T) {
	root := cmd.newRootCmd()
	root.SetArgs([]string{"diff"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestDiffCmd_MissingPlainFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", filepath.Join(tmpDir, "keys"))
	t.Setenv("ENVSEAL_STORE", filepath.Join(tmpDir, "store"))

	root := cmd.newRootCmd()
	root.SetArgs([]string{"diff", "production", filepath.Join(tmpDir, "nonexistent.env")})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing plain file")
	}
}

func TestDiffCmd_MissingSealedEnv(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", filepath.Join(tmpDir, "keys"))
	t.Setenv("ENVSEAL_STORE", filepath.Join(tmpDir, "store"))

	plainFile := filepath.Join(tmpDir, "test.env")
	if err := os.WriteFile(plainFile, []byte("FOO=bar\n"), 0600); err != nil {
		t.Fatal(err)
	}

	root := cmd.newRootCmd()
	root.SetArgs([]string{"diff", "staging", plainFile})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing sealed env")
	}
}
