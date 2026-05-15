package dotenv

import (
	"regexp"
	"strings"
)

// sensitivePattern matches common sensitive key names.
var sensitivePattern = regexp.MustCompile(
	`(?i)(password|passwd|secret|token|api_?key|private_?key|auth|credential|cert|seed|salt|hmac|signing)`,
)

// RedactedValue is the placeholder used for sensitive values.
const RedactedValue = "********"

// Redact returns a copy of the env map with sensitive values replaced by
// RedactedValue. Keys are matched case-insensitively against a set of
// well-known sensitive patterns.
func Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k) {
			out[k] = RedactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

// isSensitive reports whether the given key name looks sensitive.
func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	return sensitivePattern.MatchString(upper)
}
