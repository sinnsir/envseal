package dotenv

import (
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	src := map[string]string{"foo": "bar", "Baz": "qux"}
	r, err := Normalize(src, []NormalizeOption{NormalizeUppercaseKeys})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Output["FOO"]; !ok {
		t.Error("expected FOO in output")
	}
	if _, ok := r.Output["BAZ"]; !ok {
		t.Error("expected BAZ in output")
	}
	if len(r.Changed) != 2 {
		t.Errorf("expected 2 changed, got %d", len(r.Changed))
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	src := map[string]string{"KEY": "  hello  "}
	r, err := Normalize(src, []NormalizeOption{NormalizeTrimValues})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Output["KEY"] != "hello" {
		t.Errorf("expected 'hello', got %q", r.Output["KEY"])
	}
	if len(r.Changed) != 1 || r.Changed[0] != "KEY" {
		t.Errorf("expected KEY in changed, got %v", r.Changed)
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	src := map[string]string{"KEY": "value", "EMPTY": ""}
	r, err := Normalize(src, []NormalizeOption{NormalizeRemoveEmpty})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.Output["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if len(r.Removed) != 1 || r.Removed[0] != "EMPTY" {
		t.Errorf("expected EMPTY in removed, got %v", r.Removed)
	}
}

func TestNormalize_QuoteValues(t *testing.T) {
	src := map[string]string{"KEY": "hello world"}
	r, err := Normalize(src, []NormalizeOption{NormalizeQuoteValues})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Output["KEY"] == "hello world" {
		t.Error("expected value to be quoted")
	}
}

func TestNormalize_NoChanges(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	r, err := Normalize(src, []NormalizeOption{NormalizeTrimValues})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changes, got %v", r.Changed)
	}
	formatted := FormatNormalizeResult(r)
	if formatted != "no changes\n" {
		t.Errorf("unexpected format output: %q", formatted)
	}
}

func TestNormalize_NilSource(t *testing.T) {
	_, err := Normalize(nil, []NormalizeOption{NormalizeTrimValues})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestNormalize_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"key": "  val  "}
	_, err := Normalize(src, []NormalizeOption{NormalizeUppercaseKeys, NormalizeTrimValues})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := src["key"]; !ok {
		t.Error("source map was mutated: key 'key' missing")
	}
	if src["key"] != "  val  " {
		t.Error("source map value was mutated")
	}
}

func TestFormatNormalizeResult_WithChanges(t *testing.T) {
	r := NormalizeResult{
		Changed: []string{"FOO"},
		Removed: []string{"BAR"},
	}
	out := FormatNormalizeResult(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
