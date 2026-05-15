package dotenv

import (
	"strings"
	"testing"
)

func TestValidateSchema_AllPresent(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"SECRET_KEY":   "supersecret",
	}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Required: true},
		{Key: "SECRET_KEY", Required: true},
	}
	r := ValidateSchema(env, rules)
	if r.HasIssues() {
		t.Fatalf("expected no issues, got %+v", r)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "supersecret",
	}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Required: true},
		{Key: "SECRET_KEY", Required: true},
	}
	r := ValidateSchema(env, rules)
	if len(r.Missing) != 1 || r.Missing[0] != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL missing, got %v", r.Missing)
	}
}

func TestValidateSchema_UnknownKey(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"EXTRA_KEY":    "value",
	}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Required: true},
	}
	r := ValidateSchema(env, rules)
	if len(r.Unknown) != 1 || r.Unknown[0] != "EXTRA_KEY" {
		t.Fatalf("expected EXTRA_KEY unknown, got %v", r.Unknown)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "mysql://localhost/db",
	}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Required: true, Pattern: "postgres"},
	}
	r := ValidateSchema(env, rules)
	if len(r.Invalid) != 1 || r.Invalid[0] != "DATABASE_URL" {
		t.Fatalf("expected DATABASE_URL invalid, got %v", r.Invalid)
	}
}

func TestValidateSchema_OptionalAbsent(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
	}
	rules := []SchemaRule{
		{Key: "DATABASE_URL", Required: true},
		{Key: "LOG_LEVEL", Required: false},
	}
	r := ValidateSchema(env, rules)
	if r.HasIssues() {
		t.Fatalf("expected no issues for absent optional key, got %+v", r)
	}
}

func TestFormatSchema_NoIssues(t *testing.T) {
	r := SchemaResult{}
	out := FormatSchema(r)
	if out != "schema: ok" {
		t.Fatalf("expected 'schema: ok', got %q", out)
	}
}

func TestFormatSchema_WithIssues(t *testing.T) {
	r := SchemaResult{
		Missing: []string{"SECRET_KEY"},
		Unknown: []string{"TYPO_KEY"},
	}
	out := FormatSchema(r)
	if !strings.Contains(out, "missing required key: SECRET_KEY") {
		t.Errorf("expected missing key line, got %q", out)
	}
	if !strings.Contains(out, "unknown key: TYPO_KEY") {
		t.Errorf("expected unknown key line, got %q", out)
	}
}
