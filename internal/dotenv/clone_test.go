package dotenv

import (
	"testing"
)

func TestClone_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst, result, err := Clone(src, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst) != 3 {
		t.Errorf("expected 3 keys, got %d", len(dst))
	}
	if len(result.Keys) != 3 {
		t.Errorf("expected 3 result keys, got %d", len(result.Keys))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", result.Skipped)
	}
	// ensure it's a copy
	src["A"] = "mutated"
	if dst["A"] == "mutated" {
		t.Error("Clone should not share references with src")
	}
}

func TestClone_SubsetKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst, result, err := Clone(src, []string{"A", "C"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst) != 2 {
		t.Errorf("expected 2 keys, got %d", len(dst))
	}
	if _, ok := dst["B"]; ok {
		t.Error("key B should not be in dst")
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected no skipped, got %v", result.Skipped)
	}
}

func TestClone_MissingKeys(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst, result, err := Clone(src, []string{"A", "MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(dst) != 1 {
		t.Errorf("expected 1 key in dst, got %d", len(dst))
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "MISSING" {
		t.Errorf("expected MISSING in skipped, got %v", result.Skipped)
	}
}

func TestClone_NilSource(t *testing.T) {
	_, _, err := Clone(nil, nil)
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestFormatCloneResult_WithSkipped(t *testing.T) {
	r := CloneResult{Keys: []string{"A", "B"}, Skipped: []string{"X"}}
	out := FormatCloneResult(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
	if out == "clone: nothing to clone" {
		t.Error("unexpected empty result message")
	}
}

func TestFormatCloneResult_Empty(t *testing.T) {
	r := CloneResult{}
	out := FormatCloneResult(r)
	if out != "clone: nothing to clone" {
		t.Errorf("unexpected output: %s", out)
	}
}
