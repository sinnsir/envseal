package dotenv

import (
	"fmt"
	"strings"
)

// MaskMode controls how values are masked.
type MaskMode int

const (
	// MaskFull replaces the entire value with asterisks.
	MaskFull MaskMode = iota
	// MaskPartial reveals the first and last characters, masking the middle.
	MaskPartial
	// MaskLength replaces the value with asterisks matching the original length.
	MaskLength
)

// MaskOptions configures masking behaviour.
type MaskOptions struct {
	Mode     MaskMode
	Keys     []string // if non-empty, only mask these keys
	Exclude  []string // keys to never mask
}

// Mask returns a copy of env with sensitive values masked according to opts.
// If opts.Keys is empty, all keys considered sensitive (via isSensitive) are masked.
func Mask(env map[string]string, opts MaskOptions) map[string]string {
	excludeSet := make(map[string]bool, len(opts.Exclude))
	for _, k := range opts.Exclude {
		excludeSet[strings.ToUpper(k)] = true
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = true
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		upper := strings.ToUpper(k)
		if excludeSet[upper] {
			out[k] = v
			continue
		}
		shouldMask := len(keySet) > 0 && keySet[upper] ||
			len(keySet) == 0 && isSensitive(k)
		if shouldMask {
			out[k] = applyMask(v, opts.Mode)
		} else {
			out[k] = v
		}
	}
	return out
}

// FormatMask returns a human-readable summary of which keys were masked.
func FormatMask(original, masked map[string]string) string {
	var sb strings.Builder
	count := 0
	for k, v := range original {
		if mv, ok := masked[k]; ok && mv != v {
			count++
		}
	}
	fmt.Fprintf(&sb, "masked %d key(s)\n", count)
	for k, v := range original {
		if mv, ok := masked[k]; ok && mv != v {
			fmt.Fprintf(&sb, "  %s: %s -> %s\n", k, v, mv)
		}
	}
	return sb.String()
}

func applyMask(value string, mode MaskMode) string {
	if len(value) == 0 {
		return value
	}
	switch mode {
	case MaskPartial:
		if len(value) <= 4 {
			return strings.Repeat("*", len(value))
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	case MaskLength:
		return strings.Repeat("*", len(value))
	default: // MaskFull
		return "***"
	}
}
