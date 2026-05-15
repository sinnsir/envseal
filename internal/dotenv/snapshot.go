package dotenv

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Snapshot represents a deterministic fingerprint of a dotenv map.
type Snapshot struct {
	Hash string
	Keys []string
	Count int
}

// TakeSnapshot computes a stable SHA-256 hash over the sorted key=value pairs
// of the given map. The hash is hex-encoded and suitable for change detection.
func TakeSnapshot(m map[string]string) Snapshot {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		// Write key=value\n so each entry is unambiguous.
		fmt.Fprintf(h, "%s=%s\n", k, m[k])
	}

	return Snapshot{
		Hash:  hex.EncodeToString(h.Sum(nil)),
		Keys:  keys,
		Count: len(keys),
	}
}

// Equal reports whether two snapshots represent identical env maps.
func (s Snapshot) Equal(other Snapshot) bool {
	return s.Hash == other.Hash
}

// Short returns the first 12 hex characters of the hash, similar to a git
// short SHA, useful for display purposes.
func (s Snapshot) Short() string {
	if len(s.Hash) < 12 {
		return s.Hash
	}
	return s.Hash[:12]
}

// String implements fmt.Stringer.
func (s Snapshot) String() string {
	return fmt.Sprintf("sha256:%s (%d keys: %s)", s.Short(), s.Count, strings.Join(s.Keys, ", "))
}
