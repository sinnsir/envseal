package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// templateVarRe matches ${VAR} and $VAR style references.
var templateVarRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// TemplateResult holds the output of rendering a template map.
type TemplateResult struct {
	Rendered map[string]string
	Missing  []string
}

// Render substitutes variable references in values using the provided env map
// as the source of substitutions. References to undefined variables are left
// unexpanded and their keys are collected in TemplateResult.Missing.
func Render(vars map[string]string, env map[string]string) TemplateResult {
	result := make(map[string]string, len(vars))
	missingSet := map[string]struct{}{}

	for k, v := range vars {
		expanded := templateVarRe.ReplaceAllStringFunc(v, func(match string) string {
			name := extractVarName(match)
			if val, ok := env[name]; ok {
				return val
			}
			if val, ok := vars[name]; ok {
				return val
			}
			missingSet[name] = struct{}{}
			return match
		})
		result[k] = expanded
	}

	missing := make([]string, 0, len(missingSet))
	for k := range missingSet {
		missing = append(missing, k)
	}

	return TemplateResult{Rendered: result, Missing: missing}
}

// FormatMissing returns a human-readable summary of unresolved variables.
func FormatMissing(missing []string) string {
	if len(missing) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d unresolved variable(s):\n", len(missing)))
	for _, m := range missing {
		sb.WriteString(fmt.Sprintf("  - $%s\n", m))
	}
	return sb.String()
}

func extractVarName(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
