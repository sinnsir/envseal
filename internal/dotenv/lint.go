package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// LintIssue represents a single linting warning for a .env map.
type LintIssue struct {
	Key     string
	Message string
}

func (l LintIssue) String() string {
	return fmt.Sprintf("%s: %s", l.Key, l.Message)
}

// Lint inspects the given env map and returns a list of style/quality issues.
// It does not return errors for invalid keys (see Validate for that).
func Lint(env map[string]string) []LintIssue {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var issues []LintIssue
	for _, k := range keys {
		v := env[k]

		// Warn about keys that are not upper-case.
		if k != strings.ToUpper(k) {
			issues = append(issues, LintIssue{Key: k, Message: "key is not upper-case"})
		}

		// Warn about empty values.
		if v == "" {
			issues = append(issues, LintIssue{Key: k, Message: "value is empty"})
		}

		// Warn about values with leading or trailing whitespace.
		if strings.TrimSpace(v) != v {
			issues = append(issues, LintIssue{Key: k, Message: "value has leading or trailing whitespace"})
		}

		// Warn about values that look like they contain unexpanded variable refs.
		if strings.Contains(v, "${{") || (strings.Contains(v, "${") && strings.Contains(v, "}")) {
			issues = append(issues, LintIssue{Key: k, Message: "value may contain unexpanded variable reference"})
		}
	}
	return issues
}

// FormatLint returns a human-readable summary of lint issues.
func FormatLint(issues []LintIssue) string {
	if len(issues) == 0 {
		return "no issues found"
	}
	var sb strings.Builder
	for _, issue := range issues {
		sb.WriteString("  " + issue.String() + "\n")
	}
	return strings.TrimRight(sb.String(), "\n")
}
