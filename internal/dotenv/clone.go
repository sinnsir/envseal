package dotenv

import "fmt"

// CloneResult holds the outcome of a Clone operation.
type CloneResult struct {
	Keys    []string
	Skipped []string
}

// Clone creates a deep copy of src, optionally restricting to a subset of keys.
// If keys is empty, all keys are cloned. If a key in keys does not exist in src,
// it is recorded in Skipped.
func Clone(src map[string]string, keys []string) (map[string]string, CloneResult, error) {
	if src == nil {
		return nil, CloneResult{}, fmt.Errorf("clone: source map is nil")
	}

	dst := make(map[string]string, len(src))
	result := CloneResult{}

	if len(keys) == 0 {
		for k, v := range src {
			dst[k] = v
			result.Keys = append(result.Keys, k)
		}
		sortStrings(result.Keys)
		return dst, result, nil
	}

	for _, k := range keys {
		v, ok := src[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		dst[k] = v
		result.Keys = append(result.Keys, k)
	}

	sortStrings(result.Keys)
	sortStrings(result.Skipped)
	return dst, result, nil
}

// FormatCloneResult returns a human-readable summary of a CloneResult.
func FormatCloneResult(r CloneResult) string {
	if len(r.Keys) == 0 && len(r.Skipped) == 0 {
		return "clone: nothing to clone"
	}
	out := fmt.Sprintf("cloned %d key(s)", len(r.Keys))
	if len(r.Skipped) > 0 {
		out += fmt.Sprintf(", skipped %d missing key(s): %v", len(r.Skipped), r.Skipped)
	}
	return out
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
