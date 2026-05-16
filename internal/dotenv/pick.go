package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// PickResult holds the result of a Pick operation.
type PickResult struct {
	Picked  map[string]string
	Missing []string
}

// Pick returns a new map containing only the specified keys from src.
// If a key is not found in src, it is recorded in PickResult.Missing.
// Keys are matched case-sensitively.
func Pick(src map[string]string, keys []string) (PickResult, error) {
	if src == nil {
		return PickResult{}, fmt.Errorf("pick: source map must not be nil")
	}
	if len(keys) == 0 {
		return PickResult{}, fmt.Errorf("pick: at least one key must be specified")
	}

	picked := make(map[string]string)
	var missing []string

	for _, k := range keys {
		if v, ok := src[k]; ok {
			picked[k] = v
		} else {
			missing = append(missing, k)
		}
	}

	sort.Strings(missing)

	return PickResult{
		Picked:  picked,
		Missing: missing,
	}, nil
}

// FormatPick returns a human-readable summary of a PickResult.
func FormatPick(r PickResult) string {
	var sb strings.Builder

	keys := make([]string, 0, len(r.Picked))
	for k := range r.Picked {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(&sb, "  picked: %s\n", k)
	}
	for _, k := range r.Missing {
		fmt.Fprintf(&sb, "  missing: %s\n", k)
	}

	summary := fmt.Sprintf("picked %d key(s)", len(r.Picked))
	if len(r.Missing) > 0 {
		summary += fmt.Sprintf(", %d missing", len(r.Missing))
	}
	return summary + "\n" + sb.String()
}
