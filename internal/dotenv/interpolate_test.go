package dotenv

import (
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	res, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", res.Values["HOST"])
	}
	if len(res.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", res.Unresolved)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://${HOST}/db",
	}
	res, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["DSN"] != "postgres://localhost/db" {
		t.Errorf("unexpected DSN: %q", res.Values["DSN"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	env := map[string]string{
		"USER": "admin",
		"URL":  "http://$USER@example.com",
	}
	res, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["URL"] != "http://admin@example.com" {
		t.Errorf("unexpected URL: %q", res.Values["URL"])
	}
}

func TestInterpolate_UnresolvedVariable(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${HOST}/db",
	}
	res, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Unresolved) == 0 {
		t.Error("expected DSN to be in unresolved")
	}
}

func TestInterpolate_ChainedReferences(t *testing.T) {
	env := map[string]string{
		"A": "hello",
		"B": "${A}_world",
		"C": "${B}!",
	}
	res, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Values["C"] != "hello_world!" {
		t.Errorf("unexpected C: %q", res.Values["C"])
	}
}

func TestInterpolate_CircularReference(t *testing.T) {
	env := map[string]string{
		"A": "${B}",
		"B": "${A}",
	}
	_, err := Interpolate(env)
	if err == nil {
		t.Error("expected circular reference error, got nil")
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"BASE": "http://example.com",
		"URL":  "${BASE}/path",
	}
	original := env["URL"]
	_, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["URL"] != original {
		t.Errorf("input map was mutated: %q", env["URL"])
	}
}
