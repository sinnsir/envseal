package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// SortOrder defines how keys should be sorted.
type SortOrder int

const (
	SortAlpha      SortOrder = iota // alphabetical ascending
	SortAlphaDesc                   // alphabetical descending
	SortByLength                    // shortest key first
	SortByLengthDesc                // longest key first
)

// SortResult holds the result of a Sort operation.
type SortResult struct {
	Keys    []string
	Reorder int // number of keys that changed position
}

// Sort returns a new map with keys ordered according to the given SortOrder.
// The map itself is unordered; the sorted key list is returned in SortResult.
func Sort(src map[string]string, order SortOrder) (map[string]string, SortResult, error) {
	if src == nil {
		return nil, SortResult{}, fmt.Errorf("sort: source map is nil")
	}

	originalKeys := make([]string, 0, len(src))
	for k := range src {
		originalKeys = append(originalKeys, k)
	}
	sort.Strings(originalKeys) // stable baseline

	sortedKeys := make([]string, len(originalKeys))
	copy(sortedKeys, originalKeys)

	switch order {
	case SortAlpha:
		sort.Strings(sortedKeys)
	case SortAlphaDesc:
		sort.Sort(sort.Reverse(sort.StringSlice(sortedKeys)))
	case SortByLength:
		sort.SliceStable(sortedKeys, func(i, j int) bool {
			if len(sortedKeys[i]) == len(sortedKeys[j]) {
				return strings.ToLower(sortedKeys[i]) < strings.ToLower(sortedKeys[j])
			}
			return len(sortedKeys[i]) < len(sortedKeys[j])
		})
	case SortByLengthDesc:
		sort.SliceStable(sortedKeys, func(i, j int) bool {
			if len(sortedKeys[i]) == len(sortedKeys[j]) {
				return strings.ToLower(sortedKeys[i]) < strings.ToLower(sortedKeys[j])
			}
			return len(sortedKeys[i]) > len(sortedKeys[j])
		})
	default:
		return nil, SortResult{}, fmt.Errorf("sort: unknown order %d", order)
	}

	reorder := 0
	for i, k := range sortedKeys {
		if i < len(originalKeys) && k != originalKeys[i] {
			reorder++
		}
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}

	return out, SortResult{Keys: sortedKeys, Reorder: reorder}, nil
}

// FormatSort returns a human-readable dotenv string with keys in sorted order.
func FormatSort(src map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			continue
		}
		fmt.Fprintf(&sb, "%s=%q\n", k, v)
	}
	return sb.String()
}
