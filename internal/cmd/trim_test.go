package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestTrimCmd_NoArgs(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newTrimCmd(ks, st)
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error with no args")
	}
}

func TestTrimCmd_MissingSealedEnv(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	cmd := newTrimCmd(ks, st)
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"nonexistent"})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for missing env")
	}
}

func TestTrimCmd_DryRun(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	env := "production"
	sealEnvForTest(t, ks, st, env, map[string]string{
		"API_KEY": "  spaced  ",
		"HOST":    "clean",
	})

	out := &bytes.Buffer{}
	cmd := newTrimCmd(ks, st)
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"--leading", "--trailing", "--dry-run", env})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "API_KEY") {
		t.Errorf("expected API_KEY in dry-run output, got: %s", output)
	}
}

func TestTrimCmd_RoundTrip(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	env := "staging"
	sealEnvForTest(t, ks, st, env, map[string]string{
		"TOKEN": `"bearer-abc"`,
		"HOST":  "example.com",
	})

	out := &bytes.Buffer{}
	cmd := newTrimCmd(ks, st)
	cmd.SetOut(out)
	cmd.SetErr(out)
	cmd.SetArgs([]string{"--quotes", env})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Re-open and verify
	id, err := ks.Load(env)
	if err != nil {
		t.Fatalf("load key: %v", err)
	}
	raw, err := st.Read(env)
	if err != nil {
		t.Fatalf("read store: %v", err)
	}
	plain, err := openEnvelopeForTest(id, raw)
	if err != nil {
		t.Fatalf("open envelope: %v", err)
	}
	if !strings.Contains(string(plain), "bearer-abc") {
		t.Errorf("expected quotes stripped, got: %s", plain)
	}
}
