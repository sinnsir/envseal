package dotenv

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// GrepResult holds a single match from a grep operation.
type GrepResult struct {
	Key     string
	Value   string
	Matched string // which field matched: "key", "value", or "both"
}

// GrepOptions controls how Grep searches the env map.
type GrepOptions struct {
	Pattern     string
	SearchKeys  bool
	SearchValues bool
	IgnoreCase  bool
	Invert      bool
}

// Grep searches an env map for entries matching the given pattern.
// By default it searches both keys and values.
func Grep(env map[string]string, opts GrepOptions) ([]GrepResult, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("grep: pattern must not be empty")
	}

	pattern := opts.Pattern
	if opts.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("grep: invalid pattern %q: %w", opts.Pattern, err)
	}

	// Default: search both keys and values
	searchKeys := opts.SearchKeys || (!opts.SearchKeys && !opts.SearchValues)
	searchValues := opts.SearchValues || (!opts.SearchKeys && !opts.SearchValues)

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var results []GrepResult
	for _, k := range keys {
		v := env[k]
		keyMatch := searchKeys && re.MatchString(k)
		valMatch := searchValues && re.MatchString(v)

		matched := keyMatch || valMatch
		if opts.Invert {
			matched = !matched
		}
		if !matched {
			continue
		}

		which := ""
		if !opts.Invert {
			switch {
			case keyMatch && valMatch:
				which = "both"
			case keyMatch:
				which = "key"
			default:
				which = "value"
			}
		}

		results = append(results, GrepResult{Key: k, Value: v, Matched: which})
	}
	return results, nil
}

// FormatGrep returns a human-readable representation of grep results.
func FormatGrep(results []GrepResult) string {
	if len(results) == 0 {
		return "(no matches)"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "%s=%s\n", r.Key, r.Value)
	}
	return strings.TrimRight(sb.String(), "\n")
}
