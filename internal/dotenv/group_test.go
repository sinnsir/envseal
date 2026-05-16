package dotenv

import (
	"strings"
	"testing"
)

func TestGroup_BasicPrefixing(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "myapp",
		"SECRET": "abc",
	}
	r, err := Group(src, "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(r.Groups["DB"]))
	}
	if len(r.Groups["APP"]) != 1 {
		t.Errorf("expected 1 APP key, got %d", len(r.Groups["APP"]))
	}
	if _, ok := r.Ungrouped["SECRET"]; !ok {
		t.Error("expected SECRET in ungrouped")
	}
}

func TestGroup_DefaultSeparator(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "localhost",
		"PLAIN":   "value",
	}
	r, err := Group(src, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Groups["DB"] == nil {
		t.Error("expected DB group with default separator")
	}
}

func TestGroup_NilSource(t *testing.T) {
	_, err := Group(nil, "_")
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestGroup_AllUngrouped(t *testing.T) {
	src := map[string]string{
		"FOO": "1",
		"BAR": "2",
	}
	r, err := Group(src, "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Groups) != 0 {
		t.Errorf("expected no groups, got %d", len(r.Groups))
	}
	if len(r.Ungrouped) != 2 {
		t.Errorf("expected 2 ungrouped keys, got %d", len(r.Ungrouped))
	}
}

func TestGroup_EmptyMap(t *testing.T) {
	r, err := Group(map[string]string{}, "_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Groups) != 0 || len(r.Ungrouped) != 0 {
		t.Error("expected empty result for empty map")
	}
}

func TestFormatGroup_ContainsPrefixes(t *testing.T) {
	src := map[string]string{
		"DB_HOST": "localhost",
		"APP_ENV": "prod",
		"PLAIN":   "val",
	}
	r, _ := Group(src, "_")
	out := FormatGroup(r)
	if !strings.Contains(out, "[DB]") {
		t.Error("expected [DB] section in output")
	}
	if !strings.Contains(out, "[APP]") {
		t.Error("expected [APP] section in output")
	}
	if !strings.Contains(out, "[ungrouped]") {
		t.Error("expected [ungrouped] section in output")
	}
}
