package dotenv

import (
	"testing"
)

func TestValidate_ValidMap(t *testing.T) {
	env := map[string]string{
		"FOO":         "bar",
		"_UNDERSCORE": "ok",
		"KEY_123":     "value",
		"MixedCase":   "yes",
	}
	result := Validate(env)
	if !result.OK() {
		t.Fatalf("expected no errors, got: %s", result.Error())
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	env := map[string]string{
		"": "value",
	}
	result := Validate(env)
	if result.OK() {
		t.Fatal("expected error for empty key")
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	env := map[string]string{
		"1BAD": "value",
	}
	result := Validate(env)
	if result.OK() {
		t.Fatal("expected error for key starting with digit")
	}
}

func TestValidate_KeyWithHyphen(t *testing.T) {
	env := map[string]string{
		"BAD-KEY": "value",
	}
	result := Validate(env)
	if result.OK() {
		t.Fatal("expected error for key containing hyphen")
	}
}

func TestValidate_ValueWithNewline(t *testing.T) {
	env := map[string]string{
		"GOOD_KEY": "line1\nline2",
	}
	result := Validate(env)
	if result.OK() {
		t.Fatal("expected error for value with newline")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "GOOD_KEY" {
		t.Errorf("unexpected error key: %s", result.Errors[0].Key)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	env := map[string]string{
		"1BAD":    "value",
		"ANOTHER": "has\nnewline",
	}
	result := Validate(env)
	if result.OK() {
		t.Fatal("expected errors")
	}
	if len(result.Errors) < 2 {
		t.Fatalf("expected at least 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationError_Error(t *testing.T) {
	e := ValidationError{Key: "FOO", Message: "some issue"}
	got := e.Error()
	if got != `key "FOO": some issue` {
		t.Errorf("unexpected error string: %s", got)
	}
}
