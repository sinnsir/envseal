package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestAuditCmd_NoArgs(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newAuditCmd(ks, st)
	cmd.SetArgs([]string{})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestAuditCmd_MissingSealedEnv(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newAuditCmd(ks, st)
	cmd.SetArgs([]string{"nonexistent"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing sealed env")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention env name, got: %v", err)
	}
}

func TestAuditCmd_TextFormat(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	env := "staging"

	if err := sealEnvForTest(t, ks, st, env, map[string]string{"DB_URL": "postgres://localhost"}); err != nil {
		t.Fatalf("setup: %v", err)
	}

	cmd := newAuditCmd(ks, st)
	cmd.SetArgs([]string{env, "--format", "text"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "opened") {
		t.Errorf("expected 'opened' in output, got: %q", out)
	}
	if !strings.Contains(out, env) {
		t.Errorf("expected env name in output, got: %q", out)
	}
}

func TestAuditCmd_JSONFormat(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	env := "production"

	if err := sealEnvForTest(t, ks, st, env, map[string]string{"SECRET": "abc123"}); err != nil {
		t.Fatalf("setup: %v", err)
	}

	cmd := newAuditCmd(ks, st)
	cmd.SetArgs([]string{env, "--format", "json"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"kind"`) {
		t.Errorf("expected JSON output with 'kind' field, got: %q", out)
	}
	if !strings.Contains(out, `"opened"`) {
		t.Errorf("expected 'opened' kind in JSON, got: %q", out)
	}
}
