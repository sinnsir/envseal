package dotenv

import (
	"strings"
	"testing"
)

func TestRender_NoReferences(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := Render(vars, nil)
	if res.Rendered["FOO"] != "bar" {
		t.Errorf("expected bar, got %s", res.Rendered["FOO"])
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing, got %v", res.Missing)
	}
}

func TestRender_BraceStyle(t *testing.T) {
	vars := map[string]string{"URL": "https://${HOST}:${PORT}"}
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	res := Render(vars, env)
	if res.Rendered["URL"] != "https://localhost:8080" {
		t.Errorf("unexpected URL: %s", res.Rendered["URL"])
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing, got %v", res.Missing)
	}
}

func TestRender_DollarStyle(t *testing.T) {
	vars := map[string]string{"GREETING": "Hello $NAME"}
	env := map[string]string{"NAME": "World"}
	res := Render(vars, env)
	if res.Rendered["GREETING"] != "Hello World" {
		t.Errorf("unexpected GREETING: %s", res.Rendered["GREETING"])
	}
}

func TestRender_SelfReference(t *testing.T) {
	vars := map[string]string{"BASE": "/app", "DATA": "${BASE}/data"}
	res := Render(vars, nil)
	if res.Rendered["DATA"] != "/app/data" {
		t.Errorf("unexpected DATA: %s", res.Rendered["DATA"])
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing, got %v", res.Missing)
	}
}

func TestRender_MissingVariable(t *testing.T) {
	vars := map[string]string{"URL": "https://${UNDEFINED_HOST}"}
	res := Render(vars, nil)
	if res.Rendered["URL"] != "https://${UNDEFINED_HOST}" {
		t.Errorf("expected unexpanded value, got %s", res.Rendered["URL"])
	}
	if len(res.Missing) != 1 || res.Missing[0] != "UNDEFINED_HOST" {
		t.Errorf("expected [UNDEFINED_HOST] in missing, got %v", res.Missing)
	}
}

func TestFormatMissing_Empty(t *testing.T) {
	out := FormatMissing(nil)
	if out != "" {
		t.Errorf("expected empty string, got %q", out)
	}
}

func TestFormatMissing_WithEntries(t *testing.T) {
	out := FormatMissing([]string{"HOST", "PORT"})
	if !strings.Contains(out, "2 unresolved") {
		t.Errorf("expected count in output, got %q", out)
	}
	if !strings.Contains(out, "$HOST") || !strings.Contains(out, "$PORT") {
		t.Errorf("expected variable names in output, got %q", out)
	}
}
