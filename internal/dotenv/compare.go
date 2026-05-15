package dotenv

// CompareResult holds the result of comparing two env maps.
type CompareResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
	Same    map[string]string
}

// Compare returns a structured comparison between two env maps (old vs new).
func Compare(old, new map[string]string) CompareResult {
	result := CompareResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
		Same:    make(map[string]string),
	}

	for k, newVal := range new {
		oldVal, exists := old[k]
		if !exists {
			result.Added[k] = newVal
		} else if oldVal != newVal {
			result.Changed[k] = [2]string{oldVal, newVal}
		} else {
			result.Same[k] = newVal
		}
	}

	for k, oldVal := range old {
		if _, exists := new[k]; !exists {
			result.Removed[k] = oldVal
		}
	}

	return result
}

// HasChanges returns true if there are any differences between old and new.
func (r CompareResult) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Summary returns a short human-readable summary of the comparison.
func (r CompareResult) Summary() string {
	if !r.HasChanges() {
		return "no changes"
	}
	parts := []string{}
	if n := len(r.Added); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(r.Removed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(r.Changed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	return strings.Join(parts, ", ")
}
