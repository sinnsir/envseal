package dotenv

import "sort"

// MergeStrategy controls how conflicting keys are resolved during a merge.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing keys with incoming values.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepExisting preserves existing keys and only adds new ones.
	StrategyKeepExisting
	// StrategyError returns an error if any key conflicts.
	StrategyError
)

// ErrConflict is returned when StrategyError is used and a key conflict is found.
type ErrConflict struct {
	Key string
}

func (e *ErrConflict) Error() string {
	return "conflict: key already exists: " + e.Key
}

// MergeWithStrategy merges src into dst using the given strategy.
// It returns a new map and does not mutate the inputs.
func MergeWithStrategy(dst, src map[string]string, strategy MergeStrategy) (map[string]string, error) {
	result := make(map[string]string, len(dst))
	for k, v := range dst {
		result[k] = v
	}

	for k, v := range src {
		if _, exists := result[k]; exists {
			switch strategy {
			case StrategyOverwrite:
				result[k] = v
			case StrategyKeepExisting:
				// skip — keep dst value
			case StrategyError:
				return nil, &ErrConflict{Key: k}
			}
		} else {
			result[k] = v
		}
	}

	return result, nil
}

// Keys returns the sorted list of keys present in the map.
func Keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
