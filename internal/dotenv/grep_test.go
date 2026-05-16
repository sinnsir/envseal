package dotenv

import (
	"testing"
)

func TestGrep_MatchesValue(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"REDIS_URL":    "redis://localhost:6379",
		"APP_ENV":      "production",
	}
	results, err := Grep(env, GrepOptions{Pattern: "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestGrep_MatchesKey(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"REDIS_URL":    "redis://localhost:6379",
		"APP_ENV":      "production",
	}
	results, err := Grep(env, GrepOptions{Pattern: "URL", SearchKeys: true, SearchValues: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Matched != "key" {
			t.Errorf("expected matched=key, got %q", r.Matched)
		}
	}
}

func TestGrep_IgnoreCase(t *testing.T) {
	env := map[string]string{
		"APP_ENV": "Production",
		"DEBUG":   "false",
	}
	results, err := Grep(env, GrepOptions{Pattern: "production", IgnoreCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "APP_ENV" {
		t.Fatalf("expected APP_ENV match, got %v", results)
	}
}

func TestGrep_Invert(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"APP_ENV":      "production",
	}
	results, err := Grep(env, GrepOptions{Pattern: "localhost", Invert: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "APP_ENV" {
		t.Fatalf("expected only APP_ENV, got %v", results)
	}
}

func TestGrep_InvalidPattern(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	_, err := Grep(env, GrepOptions{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestGrep_EmptyPattern(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	_, err := Grep(env, GrepOptions{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestGrep_MatchedBoth(t *testing.T) {
	env := map[string]string{
		"URL_KEY": "url-value",
		"OTHER":   "data",
	}
	results, err := Grep(env, GrepOptions{Pattern: "url", IgnoreCase: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Matched != "both" {
		t.Errorf("expected matched=both, got %q", results[0].Matched)
	}
}

func TestFormatGrep_NoMatches(t *testing.T) {
	out := FormatGrep(nil)
	if out != "(no matches)" {
		t.Errorf("unexpected output: %q", out)
	}
}
