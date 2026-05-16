package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// HeadResult holds the output of a Head operation.
type HeadResult struct {
	Keys    []string
	Entries map[string]string
	Total   int
	Shown   int
}

// Head returns the first n key-value pairs from src, ordered alphabetically.
// If n <= 0 or n >= len(src), all entries are returned.
func Head(src map[string]string, n int) (*HeadResult, error) {
	if src == nil {
		return nil, fmt.Errorf("head: source map must not be nil")
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	total := len(keys)
	if n <= 0 || n >= total {
		n = total
	}

	selected := keys[:n]
	entries := make(map[string]string, n)
	for _, k := range selected {
		entries[k] = src[k]
	}

	return &HeadResult{
		Keys:    selected,
		Entries: entries,
		Total:   total,
		Shown:   n,
	}, nil
}

// FormatHead formats a HeadResult as a dotenv-style string.
func FormatHead(r *HeadResult) string {
	if r == nil || len(r.Keys) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, k := range r.Keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, r.Entries[k])
	}
	if r.Shown < r.Total {
		fmt.Fprintf(&sb, "# ... %d more key(s) not shown\n", r.Total-r.Shown)
	}
	return sb.String()
}
