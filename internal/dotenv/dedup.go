package dotenv

import "fmt"

// DedupStrategy controls how duplicate keys are resolved.
type DedupStrategy int

const (
	// DedupKeepFirst retains the first occurrence of a duplicate key.
	DedupKeepFirst DedupStrategy = iota
	// DedupKeepLast retains the last occurrence of a duplicate key.
	DedupKeepLast
)

// DedupResult holds the output of a Dedup operation.
type DedupResult struct {
	Out      map[string]string
	Removed  []string // keys that had duplicates removed
}

// Dedup removes duplicate keys from a slice of key=value pairs according to
// the given strategy. Unlike Parse (which silently overwrites), Dedup tracks
// which keys were duplicated and returns a DedupResult.
func Dedup(pairs []string, strategy DedupStrategy) (*DedupResult, error) {
	if pairs == nil {
		return nil, fmt.Errorf("dedup: nil input")
	}

	type entry struct {
		key   string
		value string
	}

	seen := make(map[string]int) // key -> index in ordered
	ordered := make([]entry, 0, len(pairs))
	duplicates := make(map[string]bool)

	for _, pair := range pairs {
		if pair == "" || pair[0] == '#' {
			continue
		}
		key, val, err := splitPair(pair)
		if err != nil {
			return nil, fmt.Errorf("dedup: %w", err)
		}
		if idx, exists := seen[key]; exists {
			duplicates[key] = true
			if strategy == DedupKeepLast {
				ordered[idx].value = val
			}
			// DedupKeepFirst: do nothing, keep original
		} else {
			seen[key] = len(ordered)
			ordered = append(ordered, entry{key, val})
		}
	}

	out := make(map[string]string, len(ordered))
	for _, e := range ordered {
		out[e.key] = e.value
	}

	removed := make([]string, 0, len(duplicates))
	for k := range duplicates {
		removed = append(removed, k)
	}
	sortStrings(removed)

	return &DedupResult{Out: out, Removed: removed}, nil
}

// splitPair splits a "KEY=VALUE" string into its components.
func splitPair(pair string) (string, string, error) {
	for i, ch := range pair {
		if ch == '=' {
			return pair[:i], pair[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid pair %q: missing '='" , pair)
}

// FormatDedup returns a human-readable summary of a DedupResult.
func FormatDedup(r *DedupResult) string {
	if len(r.Removed) == 0 {
		return "no duplicate keys found\n"
	}
	out := fmt.Sprintf("%d duplicate key(s) removed:\n", len(r.Removed))
	for _, k := range r.Removed {
		out += fmt.Sprintf("  - %s\n", k)
	}
	return out
}
