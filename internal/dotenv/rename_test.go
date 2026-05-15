package dotenv

import (
	"testing"
)

func TestRenameKey_Basic(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, res, err := RenameKey(src, "FOO", "FOO_NEW")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Found {
		t.Error("expected Found=true")
	}
	if _, ok := out["FOO"]; ok {
		t.Error("old key should not exist in output")
	}
	if v, ok := out["FOO_NEW"]; !ok || v != "bar" {
		t.Errorf("expected FOO_NEW=bar, got %q", v)
	}
	if out["BAZ"] != "qux" {
		t.Error("unrelated key should be preserved")
	}
}

func TestRenameKey_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	_, _, err := RenameKey(src, "FOO", "FOO_NEW")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := src["FOO"]; !ok {
		t.Error("src should not be mutated")
	}
	if _, ok := src["FOO_NEW"]; ok {
		t.Error("src should not have new key")
	}
}

func TestRenameKey_NotFound(t *testing.T) {
	src := map[string]string{"BAZ": "qux"}
	_, _, err := RenameKey(src, "MISSING", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameKey_SameKey(t *testing.T) {
	src := map[string]string{"FOO": "bar"}
	_, _, err := RenameKey(src, "FOO", "FOO")
	if err == nil {
		t.Fatal("expected error when old and new keys are the same")
	}
}

func TestRenameKey_NewKeyExists(t *testing.T) {
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	_, _, err := RenameKey(src, "FOO", "BAZ")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestRenameKey_NilSrc(t *testing.T) {
	_, _, err := RenameKey(nil, "FOO", "BAR")
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestFormatRenameResult_Found(t *testing.T) {
	r := RenameResult{OldKey: "FOO", NewKey: "FOO_NEW", Found: true}
	s := FormatRenameResult(r)
	if s == "" {
		t.Error("expected non-empty format string")
	}
}

func TestFormatRenameResult_NotFound(t *testing.T) {
	r := RenameResult{OldKey: "MISSING", NewKey: "NEW", Found: false}
	s := FormatRenameResult(r)
	if s == "" {
		t.Error("expected non-empty format string")
	}
}
