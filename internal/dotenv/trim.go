package dotenv

import (
	"fmt"
	"strings"
)

// TrimResult holds the outcome of a Trim operation.
type TrimResult struct {
	Out     map[string]string
	Trimmed []string
}

// TrimOptions controls which trimming operations are applied.
type TrimOptions struct {
	LeadingWhitespace  bool
	TrailingWhitespace bool
	Quotes             bool
	Prefix             string
	Suffix             string
}

// Trim applies the given TrimOptions to each value in src, returning a new map
// and a list of keys whose values were modified.
func Trim(src map[string]string, opts TrimOptions) (TrimResult, error) {
	if src == nil {
		return TrimResult{}, fmt.Errorf("trim: source map must not be nil")
	}
	out := make(map[string]string, len(src))
	var trimmed []string

	for k, v := range src {
		orig := v
		if opts.LeadingWhitespace {
			v = strings.TrimLeft(v, " \t")
		}
		if opts.TrailingWhitespace {
			v = strings.TrimRight(v, " \t")
		}
		if opts.Quotes {
			if len(v) >= 2 {
				if (v[0] == '"' && v[len(v)-1] == '"') ||
					(v[0] == '\'' && v[len(v)-1] == '\'') {
					v = v[1 : len(v)-1]
				}
			}
		}
		if opts.Prefix != "" {
			v = strings.TrimPrefix(v, opts.Prefix)
		}
		if opts.Suffix != "" {
			v = strings.TrimSuffix(v, opts.Suffix)
		}
		out[k] = v
		if v != orig {
			trimmed = append(trimmed, k)
		}
	}
	return TrimResult{Out: out, Trimmed: trimmed}, nil
}

// FormatTrim returns a human-readable summary of the TrimResult.
func FormatTrim(r TrimResult) string {
	if len(r.Trimmed) == 0 {
		return "no values trimmed"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d value(s) trimmed:\n", len(r.Trimmed))
	for _, k := range r.Trimmed {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return strings.TrimRight(sb.String(), "\n")
}
