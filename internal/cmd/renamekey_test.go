package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tmc/envseal/internal/dotenv"
	"github.com/tmc/envseal/internal/envelope"
)

func TestRenameKeyCmd_NoArgs(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"renamekey"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestRenameKeyCmd_MissingKey(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newRootCmd()
	cmd.SetArgs([]string{
		"--keystore", ks,
		"--store", st,
		"renamekey", "production", "FOO", "BAR",
	})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameKeyCmd_RoundTrip(t *testing.T) {
	ksDir, stDir := newTestKeystoreAndStore(t)

	// Init env
	cmd := newRootCmd()
	cmd.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "init", "staging"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("init: %v", err)
	}

	// Seal a dotenv with FOO=hello
	sealed := sealEnvForTest(t, ksDir, stDir, "staging", map[string]string{"FOO": "hello", "KEEP": "yes"})
	_ = sealed

	// Rename FOO -> FOO_RENAMED
	out := &bytes.Buffer{}
	cmd2 := newRootCmd()
	cmd2.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "renamekey", "staging", "FOO", "FOO_RENAMED"})
	cmd2.SetOut(out)
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("renamekey: %v", err)
	}
	if !strings.Contains(out.String(), "FOO_RENAMED") {
		t.Errorf("expected output to mention FOO_RENAMED, got: %q", out.String())
	}

	// Open and verify
	cmd3 := newRootCmd()
	result := &bytes.Buffer{}
	cmd3.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "export", "--format", "dotenv", "staging"})
	cmd3.SetOut(result)
	if err := cmd3.Execute(); err != nil {
		t.Fatalf("export: %v", err)
	}

	parsed, err := dotenv.Parse(result.String())
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, ok := parsed["FOO"]; ok {
		t.Error("old key FOO should not exist after rename")
	}
	if v, ok := parsed["FOO_RENAMED"]; !ok || v != "hello" {
		t.Errorf("expected FOO_RENAMED=hello, got %q", v)
	}
	if parsed["KEEP"] != "yes" {
		t.Error("unrelated key KEEP should be preserved")
	}
}

func sealEnvForTest(t *testing.T, ksDir, stDir, env string, vars map[string]string) []byte {
	t.Helper()
	_ = envelope.Seal // ensure import used
	marshaled := dotenv.Marshal(vars)

	cmd := newRootCmd()
	tf := t.TempDir() + "/test.env"
	if err := dotenv.WriteFile(tf, vars); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	cmd.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "seal", "--env", env, tf})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("seal: %v", err)
	}
	return []byte(marshaled)
}
