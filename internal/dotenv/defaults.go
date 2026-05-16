package dotenv

import "fmt"

// DefaultsResult holds the outcome of applying defaults to an env map.
type DefaultsResult struct {
	Applied map[string]string // keys that received a default value
	Skipped []string          // keys already present
}

// ApplyDefaults fills in missing keys from defaults.
// Keys already present in dst are left unchanged.
// Returns a new map and a result describing what changed.
func ApplyDefaults(dst map[string]string, defaults map[string]string) (map[string]string, DefaultsResult, error) {
	if dst == nil {
		return nil, DefaultsResult{}, fmt.Errorf("dst map must not be nil")
	}
	if defaults == nil {
		return nil, DefaultsResult{}, fmt.Errorf("defaults map must not be nil")
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	result := DefaultsResult{
		Applied: make(map[string]string),
	}

	for k, v := range defaults {
		if _, exists := out[k]; exists {
			result.Skipped = append(result.Skipped, k)
		} else {
			out[k] = v
			result.Applied[k] = v
		}
	}

	return out, result, nil
}

// FormatDefaults returns a human-readable summary of the DefaultsResult.
func FormatDefaults(r DefaultsResult) string {
	if len(r.Applied) == 0 && len(r.Skipped) == 0 {
		return "no defaults to apply\n"
	}

	var out string
	if len(r.Applied) > 0 {
		out += fmt.Sprintf("applied %d default(s):\n", len(r.Applied))
		for _, k := range Keys(r.Applied) {
			out += fmt.Sprintf("  + %s=%s\n", k, r.Applied[k])
		}
	}
	if len(r.Skipped) > 0 {
		out += fmt.Sprintf("skipped %d existing key(s):\n", len(r.Skipped))
		for _, k := range r.Skipped {
			out += fmt.Sprintf("  ~ %s\n", k)
		}
	}
	return out
}
