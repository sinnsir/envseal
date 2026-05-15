package dotenv

// PatchStrategy controls how keys are patched into a map.
type PatchStrategy int

const (
	// PatchSet adds or overwrites keys.
	PatchSet PatchStrategy = iota
	// PatchDelete removes keys.
	PatchDelete
)

// PatchOp describes a single patch operation.
type PatchOp struct {
	Key      string
	Value    string
	Strategy PatchStrategy
}

// Patch applies a slice of PatchOps to src and returns a new map.
// src is never mutated.
func Patch(src map[string]string, ops []PatchOp) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[k] = v
	}
	for _, op := range ops {
		if err := validateKey(op.Key); err != nil {
			return nil, err
		}
		switch op.Strategy {
		case PatchSet:
			out[op.Key] = op.Value
		case PatchDelete:
			delete(out, op.Key)
		}
	}
	return out, nil
}

// ParsePatchOps parses a dotenv-formatted string into a slice of PatchSet ops.
// Useful for applying a patch file expressed as a .env snippet.
func ParsePatchOps(raw string) ([]PatchOp, error) {
	m, err := Parse(raw)
	if err != nil {
		return nil, err
	}
	ops := make([]PatchOp, 0, len(m))
	for _, k := range Keys(m) {
		ops = append(ops, PatchOp{Key: k, Value: m[k], Strategy: PatchSet})
	}
	return ops, nil
}
