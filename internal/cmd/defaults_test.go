package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultsCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"defaults"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestDefaultsCmd_MissingSealedEnv(t *testing.T) {
	dir := t.TempDir()
	root := newRootCmdWithDirs(dir, dir)
	root.SetArgs([]string{"defaults", "production", "defaults.env"})
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing sealed env")
	}
}

func TestDefaultsCmd_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	ks, st := newTestKeystoreAndStore(t, dir)

	// seal an initial env
	env := map[string]string{"EXISTING": "value", "PORT": "8080"}
	if err := sealEnvForTest(t, "staging", env, ks, st); err != nil {
		t.Fatalf("seal: %v", err)
	}

	// write a defaults file
	defaultsPath := filepath.Join(dir, "defaults.env")
	content := "PORT=3000\nDEBUG=false\nLOG_LEVEL=info\n"
	if err := os.WriteFile(defaultsPath, []byte(content), 0600); err != nil {
		t.Fatalf("write defaults: %v", err)
	}

	root := newRootCmdWithDirs(dir, dir)
	root.SetArgs([]string{"defaults", "staging", defaultsPath, "--verbose"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// open and verify
	id, err := ks.Load("staging")
	if err != nil {
		t.Fatalf("load key: %v", err)
	}
	sealed, err := st.Read("staging")
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	result, err := openEnvelope(sealed, id)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if result["EXISTING"] != "value" {
		t.Errorf("EXISTING should be unchanged")
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be overwritten by default")
	}
	if result["DEBUG"] != "false" {
		t.Errorf("DEBUG should have been added from defaults")
	}
	if result["LOG_LEVEL"] != "info" {
		t.Errorf("LOG_LEVEL should have been added from defaults")
	}
}
