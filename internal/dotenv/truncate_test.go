package dotenv

import (
	"strings"
	"testing"
)

func TestTruncate_BasicShortening(t *testing.T) {
	src := map[string]string{
		"SHORT": "hi",
		"LONG":  "abcdefghij",
	}
	res, err := Truncate(src, TruncateOptions{MaxLen: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["SHORT"] != "hi" {
		t.Errorf("expected 'hi', got %q", res.Output["SHORT"])
	}
	if res.Output["LONG"] != "abcde..." {
		t.Errorf("expected 'abcde...', got %q", res.Output["LONG"])
	}
	if len(res.Truncated) != 1 || res.Truncated[0] != "LONG" {
		t.Errorf("expected [LONG] truncated, got %v", res.Truncated)
	}
}

func TestTruncate_CustomSuffix(t *testing.T) {
	src := map[string]string{"KEY": "hello world"}
	res, err := Truncate(src, TruncateOptions{MaxLen: 5, Suffix: "~~"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["KEY"] != "hello~~" {
		t.Errorf("expected 'hello~~', got %q", res.Output["KEY"])
	}
}

func TestTruncate_SpecificKeys(t *testing.T) {
	src := map[string]string{
		"A": "longvalue123",
		"B": "longvalue456",
	}
	res, err := Truncate(src, TruncateOptions{MaxLen: 4, Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["A"] != "long..." {
		t.Errorf("A should be truncated, got %q", res.Output["A"])
	}
	if res.Output["B"] != "longvalue456" {
		t.Errorf("B should be untouched, got %q", res.Output["B"])
	}
}

func TestTruncate_NilSource(t *testing.T) {
	_, err := Truncate(nil, TruncateOptions{MaxLen: 5})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestTruncate_InvalidMaxLen(t *testing.T) {
	_, err := Truncate(map[string]string{"K": "v"}, TruncateOptions{MaxLen: 0})
	if err == nil {
		t.Fatal("expected error for MaxLen=0")
	}
}

func TestTruncate_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"KEY": "abcdefgh"}
	_, err := Truncate(src, TruncateOptions{MaxLen: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["KEY"] != "abcdefgh" {
		t.Error("source map was mutated")
	}
}

func TestFormatTruncate_NoTruncated(t *testing.T) {
	r := &TruncateResult{Output: map[string]string{}, Truncated: nil}
	out := FormatTruncate(r)
	if out != "no values truncated" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatTruncate_WithTruncated(t *testing.T) {
	r := &TruncateResult{
		Output:    map[string]string{},
		Truncated: []string{"FOO", "BAR"},
	}
	out := FormatTruncate(r)
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("expected key names in output, got: %q", out)
	}
	if !strings.Contains(out, "2 key") {
		t.Errorf("expected count in output, got: %q", out)
	}
}
