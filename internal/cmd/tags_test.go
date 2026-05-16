package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
)

func TestTagsListCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"tags", "list"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing argument")
	}
}

func TestTagsFilterCmd_NoTag(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"tags", "filter", "production"})
	var buf bytes.Buffer
	root.SetErr(&buf)
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when --tag flag is missing")
	}
}

func TestTagsList_WithTags(t *testing.T) {
	ks, st := newTestKeystoreAndStore(t)

	// Build an env map that includes tag metadata.
	envMap := map[string]string{
		"DB_URL":         "postgres://localhost/db",
		"API_KEY":        "secret",
		"__TAGS__DB_URL": "tier=backend,env=prod",
	}
	id, err := ks.Load("production")
	if err != nil {
		id, err = ks.Generate("production")
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
	}
	recipient, err := id.Recipient()
	if err != nil {
		t.Fatalf("recipient: %v", err)
	}
	sealed, err := envelope.Seal(envMap, recipient)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	if err := st.Write("production", sealed); err != nil {
		t.Fatalf("write: %v", err)
	}

	root := newRootCmd()
	root.SetArgs([]string{"--keystore", ks.Dir(), "--store", st.Dir(), "tags", "list", "production"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "tier=backend") {
		t.Errorf("expected tier=backend in output, got: %s", out)
	}
}

func TestTagsFilter_ByKeyValue(t *testing.T) {
	tm := dotenv.TagMap{
		"DB_URL":  {{Key: "tier", Value: "backend"}},
		"API_KEY": {{Key: "tier", Value: "frontend"}},
	}
	env := map[string]string{"DB_URL": "postgres://", "API_KEY": "secret"}
	out := dotenv.FilterByTag(env, tm, "tier", "backend")
	if len(out) != 1 {
		t.Fatalf("expected 1 match, got %d", len(out))
	}
	if out["DB_URL"] != "postgres://" {
		t.Errorf("expected DB_URL in result")
	}
}
