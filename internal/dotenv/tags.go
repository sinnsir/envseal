package dotenv

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a key=value metadata annotation on an env var.
type Tag struct {
	Key   string
	Value string
}

// TagMap maps env var keys to their associated tags.
type TagMap map[string][]Tag

// ParseTags parses a tag string of the form "key1=val1,key2=val2".
func ParseTags(s string) ([]Tag, error) {
	if s == "" {
		return nil, nil
	}
	var tags []Tag
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid tag %q: expected key=value", part)
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if k == "" {
			return nil, fmt.Errorf("tag key must not be empty in %q", part)
		}
		tags = append(tags, Tag{Key: k, Value: v})
	}
	return tags, nil
}

// FormatTags serialises a slice of tags back to a comma-separated string.
func FormatTags(tags []Tag) string {
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = t.Key + "=" + t.Value
	}
	return strings.Join(parts, ",")
}

// FilterByTag returns a new map containing only the env vars that have a tag
// with the given key (and optionally matching value when value != "").
func FilterByTag(env map[string]string, tm TagMap, tagKey, tagValue string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		for _, t := range tm[k] {
			if t.Key == tagKey && (tagValue == "" || t.Value == tagValue) {
				out[k] = v
				break
			}
		}
	}
	return out
}

// TaggedKeys returns a sorted list of env var keys that carry at least one tag.
func TaggedKeys(tm TagMap) []string {
	keys := make([]string, 0, len(tm))
	for k := range tm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
