package dotenv

import (
	"testing"
)

func TestUnique_NoDuplicates(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	r, err := Unique(src, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Kept) != 3 {
		t.Errorf("expected 3 kept, got %d", len(r.Kept))
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(r.Removed))
	}
}

func TestUnique_RemovesAllDuplicates(t *testing.T) {
	src := map[string]string{"A": "same", "B": "same", "C": "unique"}
	r, err := Unique(src, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Kept) != 1 {
		t.Errorf("expected 1 kept, got %d", len(r.Kept))
	}
	if r.Kept["C"] != "unique" {
		t.Errorf("expected C=unique to be kept")
	}
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(r.Removed))
	}
}

func TestUnique_KeepFirst(t *testing.T) {
	src := map[string]string{"Z": "dup", "A": "dup", "M": "dup"}
	r, err := Unique(src, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Kept) != 1 {
		t.Errorf("expected 1 kept, got %d", len(r.Kept))
	}
	if _, ok := r.Kept["A"]; !ok {
		t.Errorf("expected lexicographically first key A to be kept, got %v", r.Kept)
	}
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(r.Removed))
	}
}

func TestUnique_NilSource(t *testing.T) {
	_, err := Unique(nil, false)
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestUnique_EmptyMap(t *testing.T) {
	r, err := Unique(map[string]string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Kept) != 0 || len(r.Removed) != 0 {
		t.Errorf("expected empty result for empty map")
	}
}

func TestFormatUnique_NoRemoved(t *testing.T) {
	r := UniqueResult{Kept: map[string]string{"A": "1"}, Removed: nil}
	out := FormatUnique(r)
	if out != "kept 1 key(s)" {
		t.Errorf("unexpected format: %q", out)
	}
}

func TestFormatUnique_WithRemoved(t *testing.T) {
	r := UniqueResult{
		Kept:    map[string]string{"C": "unique"},
		Removed: []string{"A", "B"},
	}
	out := FormatUnique(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
