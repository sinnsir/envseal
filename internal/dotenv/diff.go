package dotenv

import "fmt"

// DiffEntry represents a single changed key between two env maps.
type DiffEntry struct {
	Key    string
	OldVal string
	NewVal string
	Op     DiffOp
}

// DiffOp describes the type of change.
type DiffOp int

const (
	DiffAdded   DiffOp = iota // key exists in new but not old
	DiffRemoved               // key exists in old but not new
	DiffChanged               // key exists in both but value changed
)

// Diff computes the difference between two parsed env maps.
// oldMap and newMap should be produced by Parse.
func Diff(oldMap, newMap map[string]string) []DiffEntry {
	var entries []DiffEntry

	for k, newVal := range newMap {
		oldVal, exists := oldMap[k]
		if !exists {
			entries = append(entries, DiffEntry{Key: k, NewVal: newVal, Op: DiffAdded})
		} else if oldVal != newVal {
			entries = append(entries, DiffEntry{Key: k, OldVal: oldVal, NewVal: newVal, Op: DiffChanged})
		}
	}

	for k, oldVal := range oldMap {
		if _, exists := newMap[k]; !exists {
			entries = append(entries, DiffEntry{Key: k, OldVal: oldVal, Op: DiffRemoved})
		}
	}

	return entries
}

// FormatDiff returns a human-readable summary of diff entries.
func FormatDiff(entries []DiffEntry) string {
	if len(entries) == 0 {
		return "no changes"
	}
	var out string
	for _, e := range entries {
		switch e.Op {
		case DiffAdded:
			out += fmt.Sprintf("+ %s=%q\n", e.Key, e.NewVal)
		case DiffRemoved:
			out += fmt.Sprintf("- %s=%q\n", e.Key, e.OldVal)
		case DiffChanged:
			out += fmt.Sprintf("~ %s: %q -> %q\n", e.Key, e.OldVal, e.NewVal)
		}
	}
	return out
}
