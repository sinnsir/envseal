package dotenv

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationError represents a single validation issue found in a dotenv map.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) OK() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

// Validate checks that all keys in the map conform to POSIX env-var naming
// rules (start with a letter or underscore, contain only letters, digits, or
// underscores) and that no values contain unescaped newlines.
func Validate(env map[string]string) *ValidationResult {
	result := &ValidationResult{}
	for k, v := range env {
		if err := validateKey(k); err != nil {
			result.Errors = append(result.Errors, ValidationError{Key: k, Message: err.Error()})
		}
		if strings.Contains(v, "\n") {
			result.Errors = append(result.Errors, ValidationError{
				Key:     k,
				Message: "value contains unescaped newline",
			})
		}
	}
	return result
}

func validateKey(key string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	runes := []rune(key)
	first := runes[0]
	if !unicode.IsLetter(first) && first != '_' {
		return fmt.Errorf("must start with a letter or underscore, got %q", first)
	}
	for _, r := range runes[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("contains invalid character %q", r)
		}
	}
	return nil
}
