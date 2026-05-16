package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// GroupResult holds keys organized by their common prefix group.
type GroupResult struct {
	Groups map[string]map[string]string
	Ungrouped map[string]string
}

// Group partitions a dotenv map by key prefix using the given separator.
// Keys without a matching prefix are placed in Ungrouped.
// Example: Group(m, "_") groups DB_HOST, DB_PORT under "DB".
func Group(src map[string]string, sep string) (*GroupResult, error) {
	if src == nil {
		return nil, fmt.Errorf("group: source map is nil")
	}
	if sep == "" {
		sep = "_"
	}

	result := &GroupResult{
		Groups:    make(map[string]map[string]string),
		Ungrouped: make(map[string]string),
	}

	for k, v := range src {
		idx := strings.Index(k, sep)
		if idx <= 0 {
			result.Ungrouped[k] = v
			continue
		}
		prefix := k[:idx]
		if result.Groups[prefix] == nil {
			result.Groups[prefix] = make(map[string]string)
		}
		result.Groups[prefix][k] = v
	}

	return result, nil
}

// FormatGroup returns a human-readable representation of a GroupResult.
func FormatGroup(r *GroupResult) string {
	var sb strings.Builder

	prefixes := make([]string, 0, len(r.Groups))
	for p := range r.Groups {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	for _, p := range prefixes {
		fmt.Fprintf(&sb, "[%s]\n", p)
		keys := make([]string, 0, len(r.Groups[p]))
		for k := range r.Groups[p] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s=%s\n", k, r.Groups[p][k])
		}
	}

	if len(r.Ungrouped) > 0 {
		fmt.Fprintf(&sb, "[ungrouped]\n")
		keys := make([]string, 0, len(r.Ungrouped))
		for k := range r.Ungrouped {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&sb, "  %s=%s\n", k, r.Ungrouped[k])
		}
	}

	return sb.String()
}
