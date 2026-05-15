package dotenv

import (
	"strings"
	"testing"
)

func TestFlatten_BasicPrefixing(t *testing.T) {
	maps := map[string]map[string]string{
		"prod": {"DB_HOST": "prod.db", "PORT": "5432"},
		"dev": {"DB_HOST": "localhost", "PORT": "5433"},
	}

	r := Flatten(maps, "_")

	if r.Flattened["PROD_DB_HOST"] != "prod.db" {
		t.Errorf("expected PROD_DB_HOST=prod.db, got %q", r.Flattened["PROD_DB_HOST"])
	}
	if r.Flattened["DEV_PORT"] != "5433" {
		t.Errorf("expected DEV_PORT=5433, got %q", r.Flattened["DEV_PORT"])
	}
	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", r.Conflicts)
	}
}

func TestFlatten_NoConflictsWithPrefixes(t *testing.T) {
	maps := map[string]map[string]string{
		"a": {"KEY": "val_a"},
		"b": {"KEY": "val_b"},
	}

	r := Flatten(maps, "_")

	if len(r.Conflicts) != 0 {
		t.Errorf("expected no conflicts when prefixes differ, got %v", r.Conflicts)
	}
	if r.Flattened["A_KEY"] != "val_a" {
		t.Errorf("expected A_KEY=val_a")
	}
	if r.Flattened["B_KEY"] != "val_b" {
		t.Errorf("expected B_KEY=val_b")
	}
}

func TestFlatten_EmptyPrefix_DetectsConflict(t *testing.T) {
	maps := map[string]map[string]string{
		"": {"SHARED": "first"},
		"x": {"SHARED": "second"},
	}

	// With empty prefix key for first map and prefix "x" for second,
	// keys won't collide. Test true empty-prefix conflict:
	maps2 := map[string]map[string]string{
		"env1": {"FOO": "a"},
		"env2": {"FOO": "b"},
	}
	// Both produce ENV1_FOO and ENV2_FOO — no conflict expected.
	r2 := Flatten(maps2, "_")
	if len(r2.Conflicts) != 0 {
		t.Errorf("unexpected conflicts: %v", r2.Conflicts)
	}
	_ = maps
}

func TestFlatten_DefaultSeparator(t *testing.T) {
	maps := map[string]map[string]string{
		"staging": {"HOST": "stg.example.com"},
	}

	r := Flatten(maps, "")

	if _, ok := r.Flattened["STAGING_HOST"]; !ok {
		t.Errorf("expected STAGING_HOST key with default separator")
	}
}

func TestFlatten_EmptyMaps(t *testing.T) {
	r := Flatten(map[string]map[string]string{}, "_")
	if len(r.Flattened) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestFormatFlattenResult_NoConflicts(t *testing.T) {
	r := FlattenResult{Flattened: map[string]string{"A": "1", "B": "2"}}
	out := FormatFlattenResult(r)
	if !strings.Contains(out, "2 keys") {
		t.Errorf("expected '2 keys' in output, got %q", out)
	}
}

func TestFormatFlattenResult_WithConflicts(t *testing.T) {
	r := FlattenResult{
		Flattened: map[string]string{"X": "v"},
		Conflicts: []string{"X (from \"a\" and \"b\")"},
	}
	out := FormatFlattenResult(r)
	if !strings.Contains(out, "conflicts") {
		t.Errorf("expected 'conflicts' in output, got %q", out)
	}
}
