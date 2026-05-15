package dotenv

import "fmt"

// RenameResult holds the outcome of a key rename operation.
type RenameResult struct {
	OldKey string
	NewKey string
	Found  bool
}

// RenameKey renames a key in src from oldKey to newKey.
// Returns a new map with the rename applied and a RenameResult.
// Returns an error if newKey already exists or oldKey is invalid.
func RenameKey(src map[string]string, oldKey, newKey string) (map[string]string, RenameResult, error) {
	if src == nil {
		return nil, RenameResult{}, fmt.Errorf("src map is nil")
	}
	if err := validateKey(oldKey); err != nil {
		return nil, RenameResult{}, fmt.Errorf("invalid old key %q: %w", oldKey, err)
	}
	if err := validateKey(newKey); err != nil {
		return nil, RenameResult{}, fmt.Errorf("invalid new key %q: %w", newKey, err)
	}
	if oldKey == newKey {
		return nil, RenameResult{}, fmt.Errorf("old and new key are the same: %q", oldKey)
	}

	result := RenameResult{OldKey: oldKey, NewKey: newKey}

	val, found := src[oldKey]
	result.Found = found
	if !found {
		return nil, result, fmt.Errorf("key %q not found", oldKey)
	}
	if _, exists := src[newKey]; exists {
		return nil, result, fmt.Errorf("key %q already exists", newKey)
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		if k == oldKey {
			out[newKey] = val
		} else {
			out[k] = v
		}
	}
	return out, result, nil
}

// FormatRenameResult returns a human-readable summary of a RenameResult.
func FormatRenameResult(r RenameResult) string {
	if !r.Found {
		return fmt.Sprintf("key %q not found", r.OldKey)
	}
	return fmt.Sprintf("renamed %q → %q", r.OldKey, r.NewKey)
}
