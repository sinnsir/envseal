package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// SchemaRule describes a required or optional key with optional constraints.
type SchemaRule struct {
	Key      string
	Required bool
	Pattern  string // if non-empty, value must contain this substring (simple check)
}

// SchemaResult holds the outcome of a schema validation.
type SchemaResult struct {
	Missing  []string // required keys absent from the map
	Unknown  []string // keys in the map not defined in the schema
	Invalid  []string // keys whose values fail the pattern check
}

// HasIssues returns true if there are any schema violations.
func (r SchemaResult) HasIssues() bool {
	return len(r.Missing) > 0 || len(r.Unknown) > 0 || len(r.Invalid) > 0
}

// ValidateSchema checks env against the provided schema rules.
func ValidateSchema(env map[string]string, rules []SchemaRule) SchemaResult {
	ruleMap := make(map[string]SchemaRule, len(rules))
	for _, r := range rules {
		ruleMap[r.Key] = r
	}

	var result SchemaResult

	// Check required keys and pattern constraints.
	for _, r := range rules {
		val, ok := env[r.Key]
		if !ok {
			if r.Required {
				result.Missing = append(result.Missing, r.Key)
			}
			continue
		}
		if r.Pattern != "" && !strings.Contains(val, r.Pattern) {
			result.Invalid = append(result.Invalid, r.Key)
		}
	}

	// Check for unknown keys.
	for k := range env {
		if _, defined := ruleMap[k]; !defined {
			result.Unknown = append(result.Unknown, k)
		}
	}

	sort.Strings(result.Missing)
	sort.Strings(result.Unknown)
	sort.Strings(result.Invalid)
	return result
}

// FormatSchema returns a human-readable summary of schema validation results.
func FormatSchema(r SchemaResult) string {
	if !r.HasIssues() {
		return "schema: ok"
	}
	var sb strings.Builder
	for _, k := range r.Missing {
		sb.WriteString(fmt.Sprintf("missing required key: %s\n", k))
	}
	for _, k := range r.Invalid {
		sb.WriteString(fmt.Sprintf("invalid value for key: %s\n", k))
	}
	for _, k := range r.Unknown {
		sb.WriteString(fmt.Sprintf("unknown key: %s\n", k))
	}
	return strings.TrimRight(sb.String(), "\n")
}
