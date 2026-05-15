package dotenv

import (
	"testing"
)

func TestMergeWithStrategy_Overwrite(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	result, err := MergeWithStrategy(dst, src, StrategyOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "99" {
		t.Errorf("expected B=99, got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMergeWithStrategy_KeepExisting(t *testing.T) {
	dst := map[string]string{"A": "1", "B": "2"}
	src := map[string]string{"B": "99", "C": "3"}

	result, err := MergeWithStrategy(dst, src, StrategyKeepExisting)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["B"] != "2" {
		t.Errorf("expected B=2 (kept), got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestMergeWithStrategy_Error(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "2"}

	_, err := MergeWithStrategy(dst, src, StrategyError)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	conflict, ok := err.(*ErrConflict)
	if !ok {
		t.Fatalf("expected *ErrConflict, got %T", err)
	}
	if conflict.Key != "A" {
		t.Errorf("expected conflict key A, got %s", conflict.Key)
	}
}

func TestMergeWithStrategy_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"A": "1"}
	src := map[string]string{"A": "2", "B": "3"}

	_, _ = MergeWithStrategy(dst, src, StrategyOverwrite)
	if dst["A"] != "1" {
		t.Errorf("dst was mutated: A=%s", dst["A"])
	}
	if _, ok := dst["B"]; ok {
		t.Error("dst was mutated: B should not exist")
	}
}

func TestKeys_Sorted(t *testing.T) {
	m := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := Keys(m)
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}
