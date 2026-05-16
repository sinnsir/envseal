package dotenv

import (
	"testing"
)

func TestSort_Alpha(t *testing.T) {
	src := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	_, res, err := Sort(src, SortAlpha)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range res.Keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_AlphaDesc(t *testing.T) {
	src := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	_, res, err := Sort(src, SortAlphaDesc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range res.Keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_ByLength(t *testing.T) {
	src := map[string]string{"AB": "1", "ABCDE": "2", "ABC": "3"}
	_, res, err := Sort(src, SortByLength)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"AB", "ABC", "ABCDE"}
	for i, k := range res.Keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_ByLengthDesc(t *testing.T) {
	src := map[string]string{"AB": "1", "ABCDE": "2", "ABC": "3"}
	_, res, err := Sort(src, SortByLengthDesc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ABCDE", "ABC", "AB"}
	for i, k := range res.Keys {
		if k != want[i] {
			t.Errorf("pos %d: got %q, want %q", i, k, want[i])
		}
	}
}

func TestSort_NilSource(t *testing.T) {
	_, _, err := Sort(nil, SortAlpha)
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestSort_UnknownOrder(t *testing.T) {
	src := map[string]string{"A": "1"}
	_, _, err := Sort(src, SortOrder(99))
	if err == nil {
		t.Fatal("expected error for unknown order")
	}
}

func TestSort_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"B": "2", "A": "1"}
	out, _, _ := Sort(src, SortAlpha)
	out["C"] = "3"
	if _, ok := src["C"]; ok {
		t.Error("Sort mutated the source map")
	}
}

func TestFormatSort_OrderedOutput(t *testing.T) {
	src := map[string]string{"ZEBRA": "z", "APPLE": "a"}
	keys := []string{"APPLE", "ZEBRA"}
	out := FormatSort(src, keys)
	if out != "APPLE=\"a\"\nZEBRA=\"z\"\n" {
		t.Errorf("unexpected output: %q", out)
	}
}
