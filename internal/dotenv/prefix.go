package dotenv

import (
	"fmt"
	"strings"
)

// PrefixResult holds the outcome of a prefix operation.
type PrefixResult struct {
	Added   []string
	Skipped []string
}

// AddPrefix returns a new map with the given prefix added to all keys.
// Keys that already have the prefix are skipped unless force is true.
func AddPrefix(src map[string]string, prefix string, force bool) (map[string]string, PrefixResult, error) {
	if src == nil {
		return nil, PrefixResult{}, fmt.Errorf("source map is nil")
	}
	if prefix == "" {
		return nil, PrefixResult{}, fmt.Errorf("prefix must not be empty")
	}

	out := make(map[string]string, len(src))
	var result PrefixResult

	for k, v := range src {
		if !force && strings.HasPrefix(k, prefix) {
			out[k] = v
			result.Skipped = append(result.Skipped, k)
			continue
		}
		newKey := prefix + k
		out[newKey] = v
		result.Added = append(result.Added, newKey)
	}

	sortStrings(result.Added)
	sortStrings(result.Skipped)
	return out, result, nil
}

// StripPrefix returns a new map with the given prefix removed from all keys.
// Keys that do not have the prefix are skipped.
func StripPrefix(src map[string]string, prefix string) (map[string]string, PrefixResult, error) {
	if src == nil {
		return nil, PrefixResult{}, fmt.Errorf("source map is nil")
	}
	if prefix == "" {
		return nil, PrefixResult{}, fmt.Errorf("prefix must not be empty")
	}

	out := make(map[string]string, len(src))
	var result PrefixResult

	for k, v := range src {
		if !strings.HasPrefix(k, prefix) {
			out[k] = v
			result.Skipped = append(result.Skipped, k)
			continue
		}
		newKey := strings.TrimPrefix(k, prefix)
		out[newKey] = v
		result.Added = append(result.Added, newKey)
	}

	sortStrings(result.Added)
	sortStrings(result.Skipped)
	return out, result, nil
}

// FormatPrefixResult returns a human-readable summary of a prefix operation.
func FormatPrefixResult(r PrefixResult) string {
	var sb strings.Builder
	if len(r.Added) > 0 {
		sb.WriteString(fmt.Sprintf("modified %d key(s):\n", len(r.Added)))
		for _, k := range r.Added {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("skipped %d key(s) (already prefixed):\n", len(r.Skipped)))
		for _, k := range r.Skipped {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}
	return sb.String()
}
