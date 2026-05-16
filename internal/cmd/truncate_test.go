package cmd

import (
	"strings"
	"testing"
)

func TestTruncateCmd_NoArgs(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newTruncateCmd(ks, st)
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestTruncateCmd_MissingSealedEnv(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newTruncateCmd(ks, st)
	cmd.SetArgs([]string{"nonexistent"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing sealed env")
	}
}

func TestTruncateCmd_DryRun(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	env := "production"
	sealEnvForTest(t, ks, st, env, map[string]string{
		"SHORT": "hi",
		"LONG":  "this-is-a-very-long-value-that-exceeds-the-limit",
	})

	var out strings.Builder
	cmd := newTruncateCmd(ks, st)
	cmd.SetOut(&out)
	cmd.SetArgs([]string{env, "--max-len=10", "--dry-run"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "LONG") {
		t.Errorf("expected LONG in output, got: %q", out.String())
	}
}

func TestTruncateCmd_RoundTrip(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	env := "staging"
	sealEnvForTest(t, ks, st, env, map[string]string{
		"TOKEN": "supersecretlongtoken1234567890",
		"PORT":  "8080",
	})

	cmd := newTruncateCmd(ks, st)
	cmd.SetArgs([]string{env, "--max-len=10", "--suffix=~~"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// re-open and verify
	ks2, st2 := ks, st
	openCmd := newOpenCmd(ks2, st2)
	var openOut strings.Builder
	openCmd.SetOut(&openOut)
	openCmd.SetArgs([]string{env})
	if err := openCmd.Execute(); err != nil {
		t.Fatalf("open after truncate failed: %v", err)
	}
	if !strings.Contains(openOut.String(), "supersecret~~") {
		t.Errorf("expected truncated TOKEN in output, got: %q", openOut.String())
	}
	if !strings.Contains(openOut.String(), "PORT=8080") {
		t.Errorf("expected PORT unchanged, got: %q", openOut.String())
	}
}
