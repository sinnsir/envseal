package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenResult holds the output of a Flatten operation.
type FlattenResult struct {
	Flattened map[string]string
	Conflicts []string
}

// Flatten merges multiple env maps into a single map using a given prefix
// separator. Keys from later maps win unless a conflict is detected when
// strategy is "error". Prefix is prepended to all keys from each map using
// the provided separator (e.g. "_").
//
// If prefix is empty the keys are used as-is and duplicates across maps
// are treated as conflicts.
func Flatten(maps map[string]map[string]string, sep string) FlattenResult {
	if sep == "" {
		sep = "_"
	}

	out := make(map[string]string)
	var conflicts []string
	seen := make(map[string]string) // key -> source prefix

	// Sort prefixes for deterministic output.
	prefixes := make([]string, 0, len(maps))
	for p := range maps {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	for _, prefix := range prefixes {
		m := maps[prefix]
		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			var flatKey string
			if prefix == "" {
				flatKey = k
			} else {
				flatKey = strings.ToUpper(prefix) + sep + k
			}

			if src, exists := seen[flatKey]; exists {
				conflicts = append(conflicts, fmt.Sprintf("%s (from %q and %q)", flatKey, src, prefix))
			}
			seen[flatKey] = prefix
			out[flatKey] = m[k]
		}
	}

	return FlattenResult{
		Flattened: out,
		Conflicts: conflicts,
	}
}

// FormatFlattenResult returns a human-readable summary of the flatten result.
func FormatFlattenResult(r FlattenResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "flattened %d keys", len(r.Flattened))
	if len(r.Conflicts) > 0 {
		sb.WriteString("\nconflicts:\n")
		for _, c := range r.Conflicts {
			fmt.Fprintf(&sb, "  - %s\n", c)
		}
	}
	return sb.String()
}
