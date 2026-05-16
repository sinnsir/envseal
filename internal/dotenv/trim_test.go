package dotenv

import (
	"testing"
)

func TestTrim_LeadingWhitespace(t *testing.T) {
	src := map[string]string{"A": "  hello", "B": "world"}
	r, err := Trim(src, TrimOptions{LeadingWhitespace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Out["A"] != "hello" {
		t.Errorf("expected 'hello', got %q", r.Out["A"])
	}
	if len(r.Trimmed) != 1 || r.Trimmed[0] != "A" {
		t.Errorf("expected trimmed=[A], got %v", r.Trimmed)
	}
}

func TestTrim_TrailingWhitespace(t *testing.T) {
	src := map[string]string{"X": "value   ", "Y": "clean"}
	r, err := Trim(src, TrimOptions{TrailingWhitespace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Out["X"] != "value" {
		t.Errorf("expected 'value', got %q", r.Out["X"])
	}
	if r.Out["Y"] != "clean" {
		t.Errorf("Y should be unchanged")
	}
}

func TestTrim_Quotes(t *testing.T) {
	src := map[string]string{"A": `"quoted"`, "B": "'single'", "C": "bare"}
	r, err := Trim(src, TrimOptions{Quotes: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Out["A"] != "quoted" {
		t.Errorf("expected 'quoted', got %q", r.Out["A"])
	}
	if r.Out["B"] != "single" {
		t.Errorf("expected 'single', got %q", r.Out["B"])
	}
	if r.Out["C"] != "bare" {
		t.Errorf("C should be unchanged")
	}
}

func TestTrim_PrefixSuffix(t *testing.T) {
	src := map[string]string{"URL": "https://example.com/path"}
	r, err := Trim(src, TrimOptions{Prefix: "https://", Suffix: "/path"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Out["URL"] != "example.com" {
		t.Errorf("expected 'example.com', got %q", r.Out["URL"])
	}
}

func TestTrim_NilSource(t *testing.T) {
	_, err := Trim(nil, TrimOptions{})
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestTrim_NoChanges(t *testing.T) {
	src := map[string]string{"A": "clean", "B": "value"}
	r, err := Trim(src, TrimOptions{LeadingWhitespace: true, TrailingWhitespace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Trimmed) != 0 {
		t.Errorf("expected no trimmed keys, got %v", r.Trimmed)
	}
}

func TestFormatTrim_WithTrimmed(t *testing.T) {
	r := TrimResult{Out: map[string]string{"A": "v"}, Trimmed: []string{"A"}}
	out := FormatTrim(r)
	if out == "no values trimmed" {
		t.Error("expected non-empty format output")
	}
}

func TestFormatTrim_NoTrimmed(t *testing.T) {
	r := TrimResult{Out: map[string]string{}, Trimmed: nil}
	if FormatTrim(r) != "no values trimmed" {
		t.Error("expected 'no values trimmed'")
	}
}
