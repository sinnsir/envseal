package dotenv

import "fmt"

// PromoteStrategy controls how conflicting keys are handled during promotion.
type PromoteStrategy int

const (
	// PromoteSkipExisting keeps existing destination values unchanged.
	PromoteSkipExisting PromoteStrategy = iota
	// PromoteOverwrite replaces destination values with source values.
	PromoteOverwrite
)

// PromoteResult holds the outcome of a promotion operation.
type PromoteResult struct {
	Added    []string
	Skipped  []string
	Overwritten []string
}

// Promote copies keys from src into dst according to the given strategy.
// Keys present only in src are always added. Keys present in both are
// handled by the strategy. Returns a new map and a PromoteResult summary.
func Promote(src, dst map[string]string, strategy PromoteStrategy) (map[string]string, PromoteResult, error) {
	if src == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: src must not be nil")
	}
	if dst == nil {
		return nil, PromoteResult{}, fmt.Errorf("promote: dst must not be nil")
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var result PromoteResult
	for _, k := range Keys(src) {
		_, exists := out[k]
		if !exists {
			out[k] = src[k]
			result.Added = append(result.Added, k)
			continue
		}
		switch strategy {
		case PromoteOverwrite:
			out[k] = src[k]
			result.Overwritten = append(result.Overwritten, k)
		case PromoteSkipExisting:
			result.Skipped = append(result.Skipped, k)
		}
	}
	return out, result, nil
}

// FormatPromoteResult returns a human-readable summary of a PromoteResult.
func FormatPromoteResult(r PromoteResult) string {
	if len(r.Added) == 0 && len(r.Skipped) == 0 && len(r.Overwritten) == 0 {
		return "no changes\n"
	}
	var out string
	for _, k := range r.Added {
		out += fmt.Sprintf("+ %s\n", k)
	}
	for _, k := range r.Overwritten {
		out += fmt.Sprintf("~ %s\n", k)
	}
	for _, k := range r.Skipped {
		out += fmt.Sprintf("= %s (skipped)\n", k)
	}
	return out
}
