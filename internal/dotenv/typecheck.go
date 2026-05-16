package dotenv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TypeHint describes the expected type of a value.
type TypeHint string

const (
	TypeString  TypeHint = "string"
	TypeInt     TypeHint = "int"
	TypeFloat   TypeHint = "float"
	TypeBool    TypeHint = "bool"
	TypeURL     TypeHint = "url"
	TypeEmail   TypeHint = "email"
)

// TypeCheckResult holds the result for a single key.
type TypeCheckResult struct {
	Key     string
	Value   string
	Hint    TypeHint
	Valid   bool
	Message string
}

var (
	urlRe   = regexp.MustCompile(`^https?://[^\s]+$`)
	emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

// TypeCheck validates each key in src against the provided hint map.
// Keys not present in hints are skipped.
func TypeCheck(src map[string]string, hints map[string]TypeHint) []TypeCheckResult {
	var results []TypeCheckResult
	for _, k := range Keys(src) {
		hint, ok := hints[k]
		if !ok {
			continue
		}
		v := src[k]
		valid, msg := checkType(v, hint)
		results = append(results, TypeCheckResult{
			Key:     k,
			Value:   v,
			Hint:    hint,
			Valid:   valid,
			Message: msg,
		})
	}
	return results
}

func checkType(value string, hint TypeHint) (bool, string) {
	switch hint {
	case TypeInt:
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return false, fmt.Sprintf("expected integer, got %q", value)
		}
	case TypeFloat:
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return false, fmt.Sprintf("expected float, got %q", value)
		}
	case TypeBool:
		lower := strings.ToLower(value)
		valid := lower == "true" || lower == "false" || lower == "1" || lower == "0" || lower == "yes" || lower == "no"
		if !valid {
			return false, fmt.Sprintf("expected bool (true/false/1/0/yes/no), got %q", value)
		}
	case TypeURL:
		if !urlRe.MatchString(value) {
			return false, fmt.Sprintf("expected http/https URL, got %q", value)
		}
	case TypeEmail:
		if !emailRe.MatchString(value) {
			return false, fmt.Sprintf("expected email address, got %q", value)
		}
	}
	return true, ""
}

// FormatTypeCheck returns a human-readable report of type check results.
func FormatTypeCheck(results []TypeCheckResult) string {
	if len(results) == 0 {
		return "no keys checked\n"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Valid {
			fmt.Fprintf(&sb, "OK  %s (%s)\n", r.Key, r.Hint)
		} else {
			fmt.Fprintf(&sb, "ERR %s (%s): %s\n", r.Key, r.Hint, r.Message)
		}
	}
	return sb.String()
}
