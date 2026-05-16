package dotenv

import (
	"strings"
	"testing"
)

func TestTypeCheck_ValidInt(t *testing.T) {
	src := map[string]string{"PORT": "8080"}
	hints := map[string]TypeHint{"PORT": TypeInt}
	results := TypeCheck(src, hints)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Valid {
		t.Errorf("expected valid, got message: %s", results[0].Message)
	}
}

func TestTypeCheck_InvalidInt(t *testing.T) {
	src := map[string]string{"PORT": "not-a-number"}
	hints := map[string]TypeHint{"PORT": TypeInt}
	results := TypeCheck(src, hints)
	if results[0].Valid {
		t.Error("expected invalid")
	}
	if results[0].Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestTypeCheck_ValidBool(t *testing.T) {
	for _, v := range []string{"true", "false", "1", "0", "yes", "no", "TRUE", "YES"} {
		src := map[string]string{"FLAG": v}
		hints := map[string]TypeHint{"FLAG": TypeBool}
		results := TypeCheck(src, hints)
		if !results[0].Valid {
			t.Errorf("expected %q to be valid bool", v)
		}
	}
}

func TestTypeCheck_InvalidBool(t *testing.T) {
	src := map[string]string{"FLAG": "maybe"}
	hints := map[string]TypeHint{"FLAG": TypeBool}
	results := TypeCheck(src, hints)
	if results[0].Valid {
		t.Error("expected invalid bool")
	}
}

func TestTypeCheck_ValidURL(t *testing.T) {
	src := map[string]string{"API_URL": "https://example.com/api"}
	hints := map[string]TypeHint{"API_URL": TypeURL}
	results := TypeCheck(src, hints)
	if !results[0].Valid {
		t.Errorf("expected valid URL: %s", results[0].Message)
	}
}

func TestTypeCheck_InvalidURL(t *testing.T) {
	src := map[string]string{"API_URL": "not-a-url"}
	hints := map[string]TypeHint{"API_URL": TypeURL}
	results := TypeCheck(src, hints)
	if results[0].Valid {
		t.Error("expected invalid URL")
	}
}

func TestTypeCheck_ValidEmail(t *testing.T) {
	src := map[string]string{"ADMIN_EMAIL": "admin@example.com"}
	hints := map[string]TypeHint{"ADMIN_EMAIL": TypeEmail}
	results := TypeCheck(src, hints)
	if !results[0].Valid {
		t.Errorf("expected valid email: %s", results[0].Message)
	}
}

func TestTypeCheck_SkipsUnhinted(t *testing.T) {
	src := map[string]string{"FOO": "bar", "PORT": "8080"}
	hints := map[string]TypeHint{"PORT": TypeInt}
	results := TypeCheck(src, hints)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestTypeCheck_ValidFloat(t *testing.T) {
	src := map[string]string{"RATE": "3.14"}
	hints := map[string]TypeHint{"RATE": TypeFloat}
	results := TypeCheck(src, hints)
	if !results[0].Valid {
		t.Errorf("expected valid float: %s", results[0].Message)
	}
}

func TestFormatTypeCheck_OKAndErr(t *testing.T) {
	results := []TypeCheckResult{
		{Key: "PORT", Hint: TypeInt, Valid: true},
		{Key: "FLAG", Hint: TypeBool, Valid: false, Message: "expected bool"},
	}
	out := FormatTypeCheck(results)
	if !strings.Contains(out, "OK  PORT") {
		t.Errorf("expected OK line for PORT, got: %s", out)
	}
	if !strings.Contains(out, "ERR FLAG") {
		t.Errorf("expected ERR line for FLAG, got: %s", out)
	}
}

func TestFormatTypeCheck_Empty(t *testing.T) {
	out := FormatTypeCheck(nil)
	if !strings.Contains(out, "no keys checked") {
		t.Errorf("expected 'no keys checked', got: %s", out)
	}
}
