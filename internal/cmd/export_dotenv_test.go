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

func TestExportCmd_DotenvFormat(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", dir)
	t.Setenv("ENVSEAL_STORE", dir)

	ks, _ := keystore.New(dir)
	identity, err := ks.Generate("ci")
	if err != nil {
		t.Fatal(err)
	}

	original := map[string]string{"DB_URL": "postgres://localhost/mydb", "PORT": "5432"}
	plain := dotenv.Marshal(original)
	recipient, _ := identity.Recipient()
	sealed, _ := envelope.Seal(plain, recipient)

	st, _ := store.New(dir)
	if err := st.Write("ci", sealed); err != nil {
		t.Fatal(err)
	}

	out := new(bytes.Buffer)
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(new(bytes.Buffer))
	root.SetArgs([]string{"export", "ci", "--format", "dotenv"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "DB_URL") {
		t.Errorf("expected DB_URL in dotenv output, got:\n%s", got)
	}
	if !strings.Contains(got, "PORT") {
		t.Errorf("expected PORT in dotenv output, got:\n%s", got)
	}
	// dotenv format should not contain 'export' keyword
	if strings.Contains(got, "export ") {
		t.Errorf("dotenv format should not contain 'export', got:\n%s", got)
	}
}

func TestExportCmd_DefaultFormatIsShell(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE", dir)
	t.Setenv("ENVSEAL_STORE", dir)

	ks, _ := keystore.New(dir)
	identity, _ := ks.Generate("local")
	original := map[string]string{"HELLO": "world"}
	plain := dotenv.Marshal(original)
	recipient, _ := identity.Recipient()
	sealed, _ := envelope.Seal(plain, recipient)
	st, _ := store.New(dir)
	st.Write("local", sealed)

	out := new(bytes.Buffer)
	root := newRootCmd()
	root.SetOut(out)
	root.SetErr(new(bytes.Buffer))
	// no --format flag: defaults to shell
	root.SetArgs([]string{"export", "local"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "export HELLO=") {
		t.Errorf("expected shell export statement, got:\n%s", got)
	}
}
