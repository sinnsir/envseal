package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/envseal/internal/envelope"
	"github.com/yourusername/envseal/internal/keystore"
	"github.com/yourusername/envseal/internal/store"
)

func TestImportCmd_NoArgs(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"import"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestImportCmd_MissingKey(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE_DIR", tmp)
	t.Setenv("ENVSEAL_STORE_DIR", tmp)

	envFile := filepath.Join(tmp, ".env")
	_ = os.WriteFile(envFile, []byte("FOO=bar\n"), 0600)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"import", "production", envFile})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when key is missing")
	}
}

func TestImportCmd_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE_DIR", tmp)
	t.Setenv("ENVSEAL_STORE_DIR", tmp)

	ks := keystore.New(tmp)
	if _, err := ks.Generate("staging"); err != nil {
		t.Fatalf("generate key: %v", err)
	}

	cmd := newRootCmd()
	cmd.SetArgs([]string{"import", "staging", filepath.Join(tmp, "nonexistent.env")})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestImportCmd_RoundTrip(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE_DIR", tmp)
	t.Setenv("ENVSEAL_STORE_DIR", tmp)

	ks := keystore.New(tmp)
	identity, err := ks.Generate("dev")
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	envFile := filepath.Join(tmp, ".env")
	_ = os.WriteFile(envFile, []byte("APP_KEY=secret\nDEBUG=true\n"), 0600)

	buf := &strings.Builder{}
	cmd := newRootCmd()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"import", "dev", envFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("import: %v", err)
	}

	if !strings.Contains(buf.String(), "2 keys") {
		t.Errorf("expected output to mention 2 keys, got: %s", buf.String())
	}

	st := store.New(tmp)
	data, err := st.Read("dev")
	if err != nil {
		t.Fatalf("read sealed: %v", err)
	}

	env, err := envelope.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	kvs, err := envelope.Open(env, identity)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	if kvs["APP_KEY"] != "secret" {
		t.Errorf("expected APP_KEY=secret, got %q", kvs["APP_KEY"])
	}
}

func TestImportCmd_NoOverwrite(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE_DIR", tmp)
	t.Setenv("ENVSEAL_STORE_DIR", tmp)

	ks := keystore.New(tmp)
	if _, err := ks.Generate("dev"); err != nil {
		t.Fatalf("generate key: %v", err)
	}

	envFile := filepath.Join(tmp, ".env")
	_ = os.WriteFile(envFile, []byte("FOO=bar\n"), 0600)

	cmd := newRootCmd()
	cmd.SetArgs([]string{"import", "dev", envFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("first import: %v", err)
	}

	cmd2 := newRootCmd()
	cmd2.SetArgs([]string{"import", "dev", envFile})
	if err := cmd2.Execute(); err == nil {
		t.Fatal("expected error on second import without --overwrite")
	}
}
