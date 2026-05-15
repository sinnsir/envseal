package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolatePattern matches $VAR and ${VAR} style references.
var interpolatePattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// InterpolateResult holds the result of an interpolation pass.
type InterpolateResult struct {
	// Values is the interpolated map.
	Values map[string]string
	// Unresolved contains keys whose values still reference undefined variables.
	Unresolved []string
}

// Interpolate expands variable references within values using other entries in
// the same map. References that cannot be resolved are left as-is and the
// corresponding key is recorded in InterpolateResult.Unresolved.
//
// Circular references are detected and returned as an error.
func Interpolate(env map[string]string) (InterpolateResult, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	// Detect circular references via a visited set during expansion.
	var expand func(key string, seen map[string]bool) (string, error)
	expand = func(key string, seen map[string]bool) (string, error) {
		if seen[key] {
			return "", fmt.Errorf("circular reference detected for key %q", key)
		}
		val, ok := env[key]
		if !ok {
			return "$" + key, nil
		}
		seen[key] = true
		defer delete(seen, key)

		var expandErr error
		expanded := interpolatePattern.ReplaceAllStringFunc(val, func(match string) string {
			if expandErr != nil {
				return match
			}
			ref := extractRef(match)
			v, err := expand(ref, seen)
			if err != nil {
				expandErr = err
				return match
			}
			return v
		})
		if expandErr != nil {
			return "", expandErr
		}
		return expanded, nil
	}

	var unresolved []string
	for k := range env {
		expanded, err := expand(k, map[string]bool{})
		if err != nil {
			return InterpolateResult{}, err
		}
		result[k] = expanded
		if strings.Contains(expanded, "$") {
			unresolved = append(unresolved, k)
		}
	}
	return InterpolateResult{Values: result, Unresolved: unresolved}, nil
}

func extractRef(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
