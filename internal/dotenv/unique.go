package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// UniqueResult holds the result of a Unique operation.
type UniqueResult struct {
	Kept    map[string]string
	Removed []string
}

// Unique returns a new map containing only keys whose values are unique.
// When two or more keys share the same value, all are removed unless
// keepFirst is true, in which case the lexicographically first key is kept.
func Unique(src map[string]string, keepFirst bool) (UniqueResult, error) {
	if src == nil {
		return UniqueResult{}, fmt.Errorf("unique: source map is nil")
	}

	// Build a reverse index: value -> list of keys
	valueKeys := make(map[string][]string)
	for k, v := range src {
		valueKeys[v] = append(valueKeys[v], k)
	}

	kept := make(map[string]string, len(src))
	var removed []string

	for value, keys := range valueKeys {
		if len(keys) == 1 {
			kept[keys[0]] = value
			continue
		}
		sort.Strings(keys)
		if keepFirst {
			kept[keys[0]] = value
			removed = append(removed, keys[1:]...)
		} else {
			removed = append(removed, keys...)
		}
	}

	sort.Strings(removed)
	return UniqueResult{Kept: kept, Removed: removed}, nil
}

// FormatUnique returns a human-readable summary of the Unique result.
func FormatUnique(r UniqueResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "kept %d key(s)", len(r.Kept))
	if len(r.Removed) > 0 {
		fmt.Fprintf(&sb, ", removed %d duplicate(s): %s", len(r.Removed), strings.Join(r.Removed, ", "))
	}
	return sb.String()
}
