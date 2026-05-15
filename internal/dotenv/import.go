package dotenv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile reads and parses a .env file from the given path.
// It resolves the path and returns parsed key-value pairs.
func ReadFile(path string) (map[string]string, error) {
	clean := filepath.Clean(path)
	raw, err := os.ReadFile(clean)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", clean, err)
	}
	return Parse(string(raw))
}

// WriteFile marshals key-value pairs to .env format and writes to path.
func WriteFile(path string, kvs map[string]string) error {
	clean := filepath.Clean(path)
	content := Marshal(kvs)
	if err := os.WriteFile(clean, []byte(content), 0600); err != nil {
		return fmt.Errorf("writing %s: %w", clean, err)
	}
	return nil
}

// Merge combines base and override maps, with override values taking precedence.
// Returns a new map without modifying either input.
func Merge(base, override map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range override {
		out[k] = v
	}
	return out
}

// FilterKeys returns a new map containing only the keys present in the allow list.
func FilterKeys(kvs map[string]string, allow []string) map[string]string {
	set := make(map[string]bool, len(allow))
	for _, k := range allow {
		set[strings.TrimSpace(k)] = true
	}
	out := make(map[string]string)
	for k, v := range kvs {
		if set[k] {
			out[k] = v
		}
	}
	return out
}
