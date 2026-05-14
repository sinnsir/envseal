package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func TestExportCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetOut(new(bytes.Buffer))
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"export"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestExportCmd_MissingSealedEnv(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", dir)
	t.Setenv("ENVSEAL_STORE", dir)

	ks, _ := keystore.New(dir)
	_, err := ks.Generate("staging")
	if err != nil {
		t.Fatal(err)
	}

	root := newRootCmd()
	root.SetOut(new(bytes.Buffer))
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"export", "staging"})
	err = root.Execute()
	if err == nil {
		t.Fatal("expected error for missing sealed env")
	}
}

func TestExportCmd_ShellFormat(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", dir)
	t.Setenv("ENVSEAL_STORE", dir)

	ks, _ := keystore.New(dir)
	identity, err := ks.Generate("prod")
	if err != nil {
		t.Fatal(err)
	}

	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	plain := dotenv.Marshal(original)
	recipient, _ := identity.Recipient()
	sealed, _ := envelope.Seal(plain, recipient)

	st, _ := store.New(dir)
	if err := st.Write("prod", sealed); err != nil {
		t.Fatal(err)
	}

	out := new(bytes.Buffer)
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"export", "prod", "--format", "shell"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "export FOO=") {
		t.Errorf("expected shell export for FOO, got:\n%s", got)
	}
	if !strings.Contains(got, "export BAZ=") {
		t.Errorf("expected shell export for BAZ, got:\n%s", got)
	}
}

func TestExportCmd_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", dir)
	t.Setenv("ENVSEAL_STORE", dir)

	ks, _ := keystore.New(dir)
	identity, err := ks.Generate("dev")
	if err != nil {
		t.Fatal(err)
	}

	original := map[string]string{"KEY": "value"}
	plain := dotenv.Marshal(original)
	recipient, _ := identity.Recipient()
	sealed, _ := envelope.Seal(plain, recipient)

	st, _ := store.New(dir)
	st.Write("dev", sealed)

	out := new(bytes.Buffer)
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"export", "dev", "--format", "json"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, `"KEY"`) {
		t.Errorf("expected JSON with KEY, got:\n%s", got)
	}
}

func TestExportCmd_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", dir)
	t.Setenv("ENVSEAL_STORE", dir)

	ks, _ := keystore.New(dir)
	identity, _ := ks.Generate("test")
	original := map[string]string{"X": "y"}
	plain := dotenv.Marshal(original)
	recipient, _ := identity.Recipient()
	sealed, _ := envelope.Seal(plain, recipient)
	st, _ := store.New(dir)
	st.Write("test", sealed)

	root := newRootCmd()
	root.SetOut(new(bytes.Buffer))
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"export", "test", "--format", "xml"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}
