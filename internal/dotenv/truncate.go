package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// TruncateOptions controls how values are truncated.
type TruncateOptions struct {
	// MaxLen is the maximum number of characters to keep per value.
	MaxLen int
	// Suffix is appended when a value is truncated (default "...").
	Suffix string
	// Keys restricts truncation to these keys only; empty means all keys.
	Keys []string
}

// TruncateResult holds the outcome of a Truncate operation.
type TruncateResult struct {
	Output    map[string]string
	Truncated []string // keys whose values were shortened
}

// Truncate shortens values in src that exceed opts.MaxLen characters.
// It returns a new map and never mutates src.
func Truncate(src map[string]string, opts TruncateOptions) (*TruncateResult, error) {
	if src == nil {
		return nil, fmt.Errorf("truncate: source map is nil")
	}
	if opts.MaxLen <= 0 {
		return nil, fmt.Errorf("truncate: MaxLen must be greater than zero")
	}
	suffix := opts.Suffix
	if suffix == "" {
		suffix = "..."
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	out := make(map[string]string, len(src))
	var truncated []string

	for k, v := range src {
		if len(keySet) > 0 && !keySet[k] {
			out[k] = v
			continue
		}
		if len(v) > opts.MaxLen {
			out[k] = v[:opts.MaxLen] + suffix
			truncated = append(truncated, k)
		} else {
			out[k] = v
		}
	}

	sort.Strings(truncated)
	return &TruncateResult{Output: out, Truncated: truncated}, nil
}

// FormatTruncate returns a human-readable summary of a TruncateResult.
func FormatTruncate(r *TruncateResult) string {
	if len(r.Truncated) == 0 {
		return "no values truncated"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "truncated %d key(s):\n", len(r.Truncated))
	for _, k := range r.Truncated {
		fmt.Fprintf(&sb, "  %s\n", k)
	}
	return strings.TrimRight(sb.String(), "\n")
}
