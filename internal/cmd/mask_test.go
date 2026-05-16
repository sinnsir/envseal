package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestMaskCmd_NoArgs(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newRootCmd(ks, st)
	cmd.SetArgs([]string{"mask"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestMaskCmd_MissingSealedEnv(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newRootCmd(ks, st)
	cmd.SetArgs([]string{"mask", "nonexistent"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing sealed env")
	}
}

func TestMaskCmd_MasksSensitiveKeys(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	sealEnvForTest(t, ks, st, "production", map[string]string{
		"API_KEY":  "topsecret",
		"APP_HOST": "example.com",
	})

	cmd := newRootCmd(ks, st)
	cmd.SetArgs([]string{"mask", "production"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if strings.Contains(result, "topsecret") {
		t.Error("expected API_KEY value to be masked")
	}
	if !strings.Contains(result, "example.com") {
		t.Error("expected APP_HOST to be unmasked")
	}
}

func TestMaskCmd_PartialMode(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	sealEnvForTest(t, ks, st, "staging", map[string]string{
		"API_SECRET": "abcdefgh",
	})

	cmd := newRootCmd(ks, st)
	cmd.SetArgs([]string{"mask", "staging", "--mode", "partial"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if strings.Contains(result, "abcdefgh") {
		t.Error("expected API_SECRET to be partially masked")
	}
	if !strings.Contains(result, "*") {
		t.Error("expected asterisks in partial mask output")
	}
}

func TestMaskCmd_VerboseShowsSummary(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	sealEnvForTest(t, ks, st, "dev", map[string]string{
		"TOKEN": "secret123",
	})

	cmd := newRootCmd(ks, st)
	cmd.SetArgs([]string{"mask", "dev", "--verbose"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "masked") {
		t.Error("expected verbose output to mention masked keys")
	}
}
