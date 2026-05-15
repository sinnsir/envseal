package dotenv

import (
	"strings"
	"testing"
)

func TestDiff_Added(t *testing.T) {
	oldMap := map[string]string{"FOO": "bar"}
	newMap := map[string]string{"FOO": "bar", "BAZ": "qux"}

	entries := Diff(oldMap, newMap)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Op != DiffAdded || entries[0].Key != "BAZ" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	oldMap := map[string]string{"FOO": "bar", "GONE": "bye"}
	newMap := map[string]string{"FOO": "bar"}

	entries := Diff(oldMap, newMap)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Op != DiffRemoved || entries[0].Key != "GONE" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiff_Changed(t *testing.T) {
	oldMap := map[string]string{"FOO": "old"}
	newMap := map[string]string{"FOO": "new"}

	entries := Diff(oldMap, newMap)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Op != DiffChanged || e.OldVal != "old" || e.NewVal != "new" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	entries := Diff(m, m)
	if len(entries) != 0 {
		t.Errorf("expected no entries, got %d", len(entries))
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	oldMap := map[string]string{"FOO": "old", "GONE": "bye"}
	newMap := map[string]string{"FOO": "new", "NEW": "hello"}

	entries := Diff(oldMap, newMap)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	ops := map[DiffOp]int{}
	for _, e := range entries {
		ops[e.Op]++
	}
	if ops[DiffAdded] != 1 || ops[DiffRemoved] != 1 || ops[DiffChanged] != 1 {
		t.Errorf("unexpected op counts: added=%d removed=%d changed=%d",
			ops[DiffAdded], ops[DiffRemoved], ops[DiffChanged])
	}
}

func TestFormatDiff_NoChanges(t *testing.T) {
	out := FormatDiff(nil)
	if out != "no changes" {
		t.Errorf("expected 'no changes', got %q", out)
	}
}

func TestFormatDiff_Output(t *testing.T) {
	entries := []DiffEntry{
		{Key: "NEW", NewVal: "val", Op: DiffAdded},
		{Key: "OLD", OldVal: "gone", Op: DiffRemoved},
		{Key: "MOD", OldVal: "before", NewVal: "after", Op: DiffChanged},
	}
	out := FormatDiff(entries)
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("missing added line in:\n%s", out)
	}
	if !strings.Contains(out, "- OLD") {
		t.Errorf("missing removed line in:\n%s", out)
	}
	if !strings.Contains(out, "~ MOD") {
		t.Errorf("missing changed line in:\n%s", out)
	}
}
