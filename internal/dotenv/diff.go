package dotenv

import "fmt"

// DiffEntry describes a change between two Env snapshots.
type DiffEntry struct {
	Key    string
	Old    string
	New    string
	Action string // "added", "removed", "changed"
}

// Diff computes the key-level difference between two Env values.
// It returns a slice of DiffEntry describing added, removed, and changed keys.
// Values are compared semantically; comments and ordering are ignored.
func Diff(oldEnv, newEnv *Env) []DiffEntry {
	oldMap := oldEnv.Map()
	newMap := newEnv.Map()

	var diffs []DiffEntry

	for k, oldVal := range oldMap {
		newVal, exists := newMap[k]
		if !exists {
			diffs = append(diffs, DiffEntry{Key: k, Old: oldVal, Action: "removed"})
		} else if oldVal != newVal {
			diffs = append(diffs, DiffEntry{Key: k, Old: oldVal, New: newVal, Action: "changed"})
		}
	}

	for k, newVal := range newMap {
		if _, exists := oldMap[k]; !exists {
			diffs = append(diffs, DiffEntry{Key: k, New: newVal, Action: "added"})
		}
	}

	return diffs
}

// FormatDiff returns a human-readable summary of diff entries.
func FormatDiff(diffs []DiffEntry) string {
	if len(diffs) == 0 {
		return "no changes"
	}
	var out string
	for _, d := range diffs {
		switch d.Action {
		case "added":
			out += fmt.Sprintf("+ %s=%s\n", d.Key, d.New)
		case "removed":
			out += fmt.Sprintf("- %s=%s\n", d.Key, d.Old)
		case "changed":
			out += fmt.Sprintf("~ %s: %s → %s\n", d.Key, d.Old, d.New)
		}
	}
	return out
}
