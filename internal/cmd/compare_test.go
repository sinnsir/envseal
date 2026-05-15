package cmd_test

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompareCmd_NoArgs(t *testing.T) {
	root := newTestRoot(t)
	root.SetArgs([]string{"compare"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestCompareCmd_OneArg(t *testing.T) {
	root := newTestRoot(t)
	root.SetArgs([]string{"compare", "production"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for single arg")
	}
}

func TestCompareCmd_MissingKey(t *testing.T) {
	root := newTestRoot(t)
	root.SetArgs([]string{"compare", "staging", "production"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when keys are missing")
	}
}

func TestCompareCmd_RoundTrip(t *testing.T) {
	root, ksDir, stDir := newTestRootWithDirs(t)

	// init and seal staging
	root.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "init", "staging"})
	if err := root.Execute(); err != nil {
		t.Fatalf("init staging: %v", err)
	}

	stagingFile := writeTempEnvFile(t, "HELLO=world\nFOO=bar\n")
	root.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "seal", "staging", stagingFile})
	if err := root.Execute(); err != nil {
		t.Fatalf("seal staging: %v", err)
	}

	// init and seal production
	root.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "init", "production"})
	if err := root.Execute(); err != nil {
		t.Fatalf("init production: %v", err)
	}

	prodFile := writeTempEnvFile(t, "HELLO=world\nBAR=baz\n")
	root.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "seal", "production", prodFile})
	if err := root.Execute(); err != nil {
		t.Fatalf("seal production: %v", err)
	}

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetArgs([]string{"--keystore-dir", ksDir, "--store-dir", stDir, "compare", "staging", "production"})
	if err := root.Execute(); err != nil {
		t.Fatalf("compare: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "staging") || !strings.Contains(output, "production") {
		t.Errorf("expected env names in output, got: %s", output)
	}
	if !strings.Contains(output, "Summary:") {
		t.Errorf("expected Summary in output, got: %s", output)
	}
}
