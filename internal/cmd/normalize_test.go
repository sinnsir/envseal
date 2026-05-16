package cmd

import (
	"bytes"
	"testing"
)

func TestNormalizeCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"normalize"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestNormalizeCmd_MissingSealedEnv(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)
	root := newRootCmd()
	root.SetArgs([]string{
		"--keystore", ks,
		"--store", st,
		"normalize", "production",
	})
	var buf bytes.Buffer
	root.SetErr(&buf)
	if err := root.Execute(); err == nil {
		t.Error("expected error for missing sealed env")
	}
}

func TestNormalizeCmd_DryRun(t *testing.T) {
	ksDir, stDir := newTestKeystoreAndStore(t)

	root := newRootCmd()
	root.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "init", "staging"})
	if err := root.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	root2 := newRootCmd()
	root2.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "seal", "staging", "--input", "-"})
	var sealBuf bytes.Buffer
	root2.SetIn(bytes.NewBufferString("key=  hello  \nEMPTY=\n"))
	root2.SetOut(&sealBuf)
	if err := root2.Execute(); err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	root3 := newRootCmd()
	root3.SetArgs([]string{
		"--keystore", ksDir,
		"--store", stDir,
		"normalize", "staging",
		"--trim",
		"--remove-empty",
		"--dry-run",
	})
	var out bytes.Buffer
	root3.SetOut(&out)
	if err := root3.Execute(); err != nil {
		t.Fatalf("normalize dry-run failed: %v", err)
	}
}

func TestNormalizeCmd_RoundTrip(t *testing.T) {
	ksDir, stDir := newTestKeystoreAndStore(t)

	root := newRootCmd()
	root.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "init", "dev"})
	if err := root.Execute(); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	root2 := newRootCmd()
	root2.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "seal", "dev", "--input", "-"})
	root2.SetIn(bytes.NewBufferString("lower_key=value\n"))
	var sealOut bytes.Buffer
	root2.SetOut(&sealOut)
	if err := root2.Execute(); err != nil {
		t.Fatalf("seal failed: %v", err)
	}

	root3 := newRootCmd()
	root3.SetArgs([]string{
		"--keystore", ksDir,
		"--store", stDir,
		"normalize", "dev",
		"--uppercase",
	})
	var normOut bytes.Buffer
	root3.SetOut(&normOut)
	if err := root3.Execute(); err != nil {
		t.Fatalf("normalize failed: %v", err)
	}

	root4 := newRootCmd()
	root4.SetArgs([]string{"--keystore", ksDir, "--store", stDir, "export", "dev", "--format", "dotenv"})
	var exportOut bytes.Buffer
	root4.SetOut(&exportOut)
	if err := root4.Execute(); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	output := exportOut.String()
	if len(output) == 0 {
		t.Error("expected non-empty export output after normalize")
	}
}
