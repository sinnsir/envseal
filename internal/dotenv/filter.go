package dotenv

import (
	"regexp"
	"strings"
)

// FilterOptions controls how Filter selects keys from a map.
type FilterOptions struct {
	// Prefix retains only keys with the given prefix.
	Prefix string
	// Suffix retains only keys with the given suffix.
	Suffix string
	// Pattern retains only keys matching the given regular expression.
	Pattern string
	// Invert negates the filter, retaining keys that do NOT match.
	Invert bool
}

// Filter returns a new map containing only the entries whose keys satisfy
// the criteria defined in opts. All criteria are ANDed together.
func Filter(src map[string]string, opts FilterOptions) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
	}

	out := make(map[string]string, len(src))
	for k, v := range src {
		matches := matchesFilter(k, opts, re)
		if opts.Invert {
			matches = !matches
		}
		if matches {
			out[k] = v
		}
	}
	return out, nil
}

func matchesFilter(key string, opts FilterOptions, re *regexp.Regexp) bool {
	if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
		return false
	}
	if opts.Suffix != "" && !strings.HasSuffix(key, opts.Suffix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}
