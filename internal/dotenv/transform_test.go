package dotenv

import (
	"testing"
)

func TestTransform_Uppercase(t *testing.T) {
	src := map[string]string{"KEY": "hello", "OTHER": "world"}
	got, err := Transform(src, []string{"uppercase"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "HELLO" || got["OTHER"] != "WORLD" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestTransform_Lowercase(t *testing.T) {
	src := map[string]string{"KEY": "HELLO"}
	got, err := Transform(src, []string{"lowercase"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "hello" {
		t.Errorf("got %q, want %q", got["KEY"], "hello")
	}
}

func TestTransform_Trim(t *testing.T) {
	src := map[string]string{"KEY": "  spaced  "}
	got, err := Transform(src, []string{"trim"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "spaced" {
		t.Errorf("got %q, want %q", got["KEY"], "spaced")
	}
}

func TestTransform_StripQuotes(t *testing.T) {
	src := map[string]string{"A": `"quoted"`, "B": "'single'", "C": "bare"}
	got, err := Transform(src, []string{"strip_quotes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["A"] != "quoted" {
		t.Errorf("A: got %q, want %q", got["A"], "quoted")
	}
	if got["B"] != "single" {
		t.Errorf("B: got %q, want %q", got["B"], "single")
	}
	if got["C"] != "bare" {
		t.Errorf("C: got %q, want %q", got["C"], "bare")
	}
}

func TestTransform_ChainedOps(t *testing.T) {
	src := map[string]string{"KEY": "  Hello World  "}
	got, err := Transform(src, []string{"trim", "uppercase"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "HELLO WORLD" {
		t.Errorf("got %q, want %q", got["KEY"], "HELLO WORLD")
	}
}

func TestTransform_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	_, err := Transform(src, []string{"uppercase"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src["KEY"] != "value" {
		t.Error("Transform mutated the source map")
	}
}

func TestTransform_UnknownOp(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	_, err := Transform(src, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown transform, got nil")
	}
}

func TestTransformKeys_ReturnsSorted(t *testing.T) {
	keys := TransformKeys()
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("TransformKeys not sorted: %v", keys)
			break
		}
	}
	if len(keys) == 0 {
		t.Error("expected at least one transform key")
	}
}
