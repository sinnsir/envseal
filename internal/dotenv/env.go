package dotenv

import (
	"fmt"
	"os"
	"sort"
)

// FromEnv returns a map of environment variables from the current process
// filtered to only include keys present in the provided set.
// If keys is nil or empty, all environment variables are returned.
func FromEnv(keys []string) map[string]string {
	all := os.Environ()
	result := make(map[string]string, len(all))

	for _, entry := range all {
		for i := 0; i < len(entry); i++ {
			if entry[i] == '=' {
				result[entry[:i]] = entry[i+1:]
				break
			}
		}
	}

	if len(keys) == 0 {
		return result
	}

	filtered := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := result[k]; ok {
			filtered[k] = v
		}
	}
	return filtered
}

// ToEnv converts a map to a slice of KEY=VALUE strings suitable for
// use with os/exec Env fields. Keys are sorted for determinism.
func ToEnv(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	env := make([]string, 0, len(m))
	for _, k := range keys {
		env = append(env, fmt.Sprintf("%s=%s", k, m[k]))
	}
	return env
}

// ApplyToProcess sets each key-value pair in m as an environment variable
// in the current process. Returns the first error encountered, if any.
func ApplyToProcess(m map[string]string) error {
	for k, v := range m {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("setenv %q: %w", k, err)
		}
	}
	return nil
}
