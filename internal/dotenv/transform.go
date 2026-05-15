package dotenv

import (
	"fmt"
	"strings"
)

// TransformFn is a function that transforms a single env value.
type TransformFn func(key, value string) (string, error)

// TransformOp describes a named transformation to apply.
type TransformOp struct {
	Name string
	Fn   TransformFn
}

// BuiltinTransforms contains the available named transforms.
var BuiltinTransforms = map[string]TransformFn{
	"uppercase": func(_, v string) (string, error) { return strings.ToUpper(v), nil },
	"lowercase": func(_, v string) (string, error) { return strings.ToLower(v), nil },
	"trim":      func(_, v string) (string, error) { return strings.TrimSpace(v), nil },
	"strip_quotes": func(_, v string) (string, error) {
		v = strings.TrimSpace(v)
		if len(v) >= 2 {
			if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
				return v[1 : len(v)-1], nil
			}
		}
		return v, nil
	},
}

// Transform applies one or more named transform operations to every value in src.
// It returns a new map without mutating src.
func Transform(src map[string]string, ops []string) (map[string]string, error) {
	fns := make([]TransformFn, 0, len(ops))
	for _, name := range ops {
		fn, ok := BuiltinTransforms[name]
		if !ok {
			return nil, fmt.Errorf("unknown transform %q", name)
		}
		fns = append(fns, fn)
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		cur := v
		for _, fn := range fns {
			var err error
			cur, err = fn(k, cur)
			if err != nil {
				return nil, fmt.Errorf("transform key %q: %w", k, err)
			}
		}
		out[k] = cur
	}
	return out, nil
}

// TransformKeys returns the sorted list of available built-in transform names.
func TransformKeys() []string {
	keys := make([]string, 0, len(BuiltinTransforms))
	for k := range BuiltinTransforms {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
