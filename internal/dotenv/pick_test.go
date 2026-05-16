package dotenv

import (
	"strings"
	"testing"
)

func TestPick_BasicSubset(t *testing.T) {
	src := map[string]string{
		"FOO": "foo",
		"BAR": "bar",
		"BAZ": "baz",
	}
	r, err := Pick(src, []string{"FOO", "BAZ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Picked) != 2 {
		t.Errorf("expected 2 picked, got %d", len(r.Picked))
	}
	if r.Picked["FOO"] != "foo" || r.Picked["BAZ"] != "baz" {
		t.Errorf("unexpected picked values: %v", r.Picked)
	}
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing, got %v", r.Missing)
	}
}

func TestPick_MissingKeys(t *testing.T) {
	src := map[string]string{"FOO": "foo"}
	r, err := Pick(src, []string{"FOO", "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Picked) != 1 {
		t.Errorf("expected 1 picked, got %d", len(r.Picked))
	}
	if len(r.Missing) != 1 || r.Missing[0] != "MISSING" {
		t.Errorf("expected [MISSING], got %v", r.Missing)
	}
}

func TestPick_AllMissing(t *testing.T) {
	src := map[string]string{"FOO": "foo"}
	r, err := Pick(src, []string{"A", "B"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Picked) != 0 {
		t.Errorf("expected 0 picked, got %d", len(r.Picked))
	}
	if len(r.Missing) != 2 {
		t.Errorf("expected 2 missing, got %v", r.Missing)
	}
}

func TestPick_NilSource(t *testing.T) {
	_, err := Pick(nil, []string{"FOO"})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestPick_NoKeys(t *testing.T) {
	src := map[string]string{"FOO": "foo"}
	_, err := Pick(src, []string{})
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestPick_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"FOO": "foo", "BAR": "bar"}
	r, _ := Pick(src, []string{"FOO"})
	r.Picked["FOO"] = "mutated"
	if src["FOO"] != "foo" {
		t.Error("Pick mutated the source map")
	}
}

func TestFormatPick_ContainsPickedAndMissing(t *testing.T) {
	r := PickResult{
		Picked:  map[string]string{"FOO": "foo"},
		Missing: []string{"BAR"},
	}
	out := FormatPick(r)
	if !strings.Contains(out, "picked: FOO") {
		t.Errorf("expected 'picked: FOO' in output, got: %s", out)
	}
	if !strings.Contains(out, "missing: BAR") {
		t.Errorf("expected 'missing: BAR' in output, got: %s", out)
	}
	if !strings.Contains(out, "1 missing") {
		t.Errorf("expected '1 missing' in summary, got: %s", out)
	}
}
