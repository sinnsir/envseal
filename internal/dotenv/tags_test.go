package dotenv

import (
	"testing"
)

func TestParseTags_Basic(t *testing.T) {
	tags, err := ParseTags("env=prod,tier=backend")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Key != "env" || tags[0].Value != "prod" {
		t.Errorf("unexpected tag[0]: %+v", tags[0])
	}
	if tags[1].Key != "tier" || tags[1].Value != "backend" {
		t.Errorf("unexpected tag[1]: %+v", tags[1])
	}
}

func TestParseTags_Empty(t *testing.T) {
	tags, err := ParseTags("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected empty tags, got %v", tags)
	}
}

func TestParseTags_InvalidNoEquals(t *testing.T) {
	_, err := ParseTags("badtag")
	if err == nil {
		t.Fatal("expected error for tag without '='")
	}
}

func TestParseTags_EmptyKey(t *testing.T) {
	_, err := ParseTags("=value")
	if err == nil {
		t.Fatal("expected error for empty tag key")
	}
}

func TestFormatTags_RoundTrip(t *testing.T) {
	input := "env=prod,tier=backend"
	tags, err := ParseTags(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := FormatTags(tags)
	if got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestFilterByTag_MatchKey(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2", "BAZ": "3"}
	tm := TagMap{
		"FOO": {{Key: "env", Value: "prod"}},
		"BAZ": {{Key: "env", Value: "staging"}},
	}
	out := FilterByTag(env, tm, "env", "")
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	if _, ok := out["BAR"]; ok {
		t.Error("BAR should not be in result")
	}
}

func TestFilterByTag_MatchKeyAndValue(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAZ": "3"}
	tm := TagMap{
		"FOO": {{Key: "env", Value: "prod"}},
		"BAZ": {{Key: "env", Value: "staging"}},
	}
	out := FilterByTag(env, tm, "env", "prod")
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out["FOO"] != "1" {
		t.Errorf("expected FOO=1, got %v", out["FOO"])
	}
}

func TestTaggedKeys_Sorted(t *testing.T) {
	tm := TagMap{
		"ZZZ": {{Key: "x", Value: "1"}},
		"AAA": {{Key: "x", Value: "2"}},
		"MMM": {{Key: "x", Value: "3"}},
	}
	keys := TaggedKeys(tm)
	if keys[0] != "AAA" || keys[1] != "MMM" || keys[2] != "ZZZ" {
		t.Errorf("unexpected order: %v", keys)
	}
}
