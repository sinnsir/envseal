package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// NormalizeOption controls how normalization is applied.
type NormalizeOption int

const (
	NormalizeUppercaseKeys NormalizeOption = iota
	NormalizeTrimValues
	NormalizeRemoveEmpty
	NormalizeQuoteValues
)

// NormalizeResult holds the outcome of a normalization run.
type NormalizeResult struct {
	Output   map[string]string
	Changed  []string // keys that were modified
	Removed  []string // keys that were dropped
}

// Normalize applies a set of normalization options to an env map and returns
// a new map along with a result describing what changed.
func Normalize(src map[string]string, opts []NormalizeOption) (NormalizeResult, error) {
	if src == nil {
		return NormalizeResult{}, fmt.Errorf("normalize: source map is nil")
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	optSet := make(map[NormalizeOption]bool, len(opts))
	for _, o := range opts {
		optSet[o] = true
	}

	var changed, removed []string

	// Uppercase keys first (may rename keys)
	if optSet[NormalizeUppercaseKeys] {
		for k, v := range out {
			upper := strings.ToUpper(k)
			if upper != k {
				delete(out, k)
				out[upper] = v
				changed = append(changed, upper)
			}
		}
	}

	for k, v := range out {
		original := v

		if optSet[NormalizeTrimValues] {
			v = strings.TrimSpace(v)
		}

		if optSet[NormalizeRemoveEmpty] && v == "" {
			delete(out, k)
			removed = append(removed, k)
			continue
		}

		if optSet[NormalizeQuoteValues] && strings.ContainsAny(v, " \t") {
			v = fmt.Sprintf("%q", v)
		}

		if v != original {
			out[k] = v
			if !contains(changed, k) {
				changed = append(changed, k)
			}
		}
	}

	sort.Strings(changed)
	sort.Strings(removed)

	return NormalizeResult{Output: out, Changed: changed, Removed: removed}, nil
}

// FormatNormalizeResult returns a human-readable summary of the normalization.
func FormatNormalizeResult(r NormalizeResult) string {
	var sb strings.Builder
	if len(r.Changed) == 0 && len(r.Removed) == 0 {
		sb.WriteString("no changes\n")
		return sb.String()
	}
	for _, k := range r.Changed {
		fmt.Fprintf(&sb, "~ %s\n", k)
	}
	for _, k := range r.Removed {
		fmt.Fprintf(&sb, "- %s\n", k)
	}
	return sb.String()
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
