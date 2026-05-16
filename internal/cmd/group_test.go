package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestGroupCmd_NoArgs(t *testing.T) {
	cmd := newGroupCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error with no args")
	}
}

func TestGroupCmd_MissingSealedEnv(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	root := newRootCmd()
	root.SetArgs([]string{
		"--keystore-dir", ks,
		"--store-dir", st,
		"group", "nonexistent",
	})
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})
	err := root.Execute()
	if err == nil {
		t.Error("expected error for missing sealed env")
	}
}

func TestGroupCmd_RoundTrip(t *testing.T) {
	ksDir, stDir := newTestKeystoreAndStore(t)

	// Init key and seal env
	env := "production"
	plain := "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\nSECRET=abc\n"
	sealEnvForTest(t, ksDir, stDir, env, plain)

	root := newRootCmd()
	root.SetArgs([]string{
		"--keystore-dir", ksDir,
		"--store-dir", stDir,
		"group", env,
	})
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outStr := out.String()
	if !strings.Contains(outStr, "[DB]") {
		t.Errorf("expected [DB] group in output, got:\n%s", outStr)
	}
	if !strings.Contains(outStr, "[APP]") {
		t.Errorf("expected [APP] group in output, got:\n%s", outStr)
	}
	if !strings.Contains(outStr, "[ungrouped]") {
		t.Errorf("expected [ungrouped] section in output, got:\n%s", outStr)
	}
}

func TestGroupCmd_ListFormat(t *testing.T) {
	ksDir, stDir := newTestKeystoreAndStore(t)

	env := "staging"
	plain := "DB_HOST=localhost\nDB_PORT=5432\n"
	sealEnvForTest(t, ksDir, stDir, env, plain)

	root := newRootCmd()
	root.SetArgs([]string{
		"--keystore-dir", ksDir,
		"--store-dir", stDir,
		"group", env, "--format", "list",
	})
	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outStr := out.String()
	if !strings.Contains(outStr, "DB") {
		t.Errorf("expected DB prefix in list output, got:\n%s", outStr)
	}
}
