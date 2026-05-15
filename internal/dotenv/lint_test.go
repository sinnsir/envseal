package dotenv

import (
	"strings"
	"testing"
)

func TestLint_NoIssues(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"PORT":         "8080",
	}
	issues := Lint(env)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{"database_url": "postgres://localhost"}
	issues := Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "upper-case") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	issues := Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "empty") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_WhitespaceValue(t *testing.T) {
	env := map[string]string{"SECRET": "  value  "}
	issues := Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "whitespace") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_UnexpandedVariable(t *testing.T) {
	env := map[string]string{"ENDPOINT": "https://api.${HOST}/v1"}
	issues := Lint(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "unexpanded") {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"bad_key": "",
		"GOOD_KEY": "ok",
	}
	issues := Lint(env)
	// bad_key triggers lowercase + empty value = 2 issues
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issues), issues)
	}
}

func TestFormatLint_NoIssues(t *testing.T) {
	out := FormatLint(nil)
	if out != "no issues found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatLint_WithIssues(t *testing.T) {
	issues := []LintIssue{
		{Key: "FOO", Message: "value is empty"},
	}
	out := FormatLint(issues)
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected key in output, got: %q", out)
	}
	if !strings.Contains(out, "value is empty") {
		t.Errorf("expected message in output, got: %q", out)
	}
}
