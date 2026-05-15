package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// Summary holds statistics about a dotenv map.
type Summary struct {
	Total     int
	Empty     int
	Sensitive int
	Keys      []string
}

// Summarize returns a Summary of the given env map.
func Summarize(env map[string]string) Summary {
	keys := make([]string, 0, len(env))
	var empty, sensitive int

	for k, v := range env {
		keys = append(keys, k)
		if strings.TrimSpace(v) == "" {
			empty++
		}
		if isSensitive(k) {
			sensitive++
		}
	}
	sort.Strings(keys)

	return Summary{
		Total:     len(env),
		Empty:     empty,
		Sensitive: sensitive,
		Keys:      keys,
	}
}

// FormatSummary returns a human-readable summary string.
func FormatSummary(s Summary) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Total keys : %d\n", s.Total)
	fmt.Fprintf(&sb, "Empty values: %d\n", s.Empty)
	fmt.Fprintf(&sb, "Sensitive   : %d\n", s.Sensitive)
	if len(s.Keys) > 0 {
		sb.WriteString("Keys:\n")
		for _, k := range s.Keys {
			fmt.Fprintf(&sb, "  - %s\n", k)
		}
	}
	return sb.String()
}
