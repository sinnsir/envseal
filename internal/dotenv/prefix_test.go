package dotenv

import (
	"testing"
)

func TestAddPrefix_Basic(t *testing.T) {
	src := map[string]string{"FOO": "1", "BAR": "2"}
	out, result, err := AddPrefix(src, "APP_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_FOO"] != "1" || out["APP_BAR"] != "2" {
		t.Errorf("expected prefixed keys, got %v", out)
	}
	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(result.Skipped))
	}
}

func TestAddPrefix_SkipsAlreadyPrefixed(t *testing.T) {
	src := map[string]string{"APP_FOO": "1", "BAR": "2"}
	out, result, err := AddPrefix(src, "APP_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_FOO"] != "1" || out["APP_BAR"] != "2" {
		t.Errorf("unexpected output: %v", out)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "APP_FOO" {
		t.Errorf("expected APP_FOO skipped, got %v", result.Skipped)
	}
}

func TestAddPrefix_Force(t *testing.T) {
	src := map[string]string{"APP_FOO": "1"}
	out, result, err := AddPrefix(src, "APP_", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_APP_FOO"] != "1" {
		t.Errorf("expected double-prefixed key, got %v", out)
	}
	if len(result.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Added))
	}
}

func TestAddPrefix_NilSource(t *testing.T) {
	_, _, err := AddPrefix(nil, "APP_", false)
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestAddPrefix_EmptyPrefix(t *testing.T) {
	_, _, err := AddPrefix(map[string]string{"K": "v"}, "", false)
	if err == nil {
		t.Error("expected error for empty prefix")
	}
}

func TestStripPrefix_Basic(t *testing.T) {
	src := map[string]string{"APP_FOO": "1", "APP_BAR": "2", "OTHER": "3"}
	out, result, err := StripPrefix(src, "APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "1" || out["BAR"] != "2" {
		t.Errorf("expected stripped keys, got %v", out)
	}
	if out["OTHER"] != "3" {
		t.Errorf("expected OTHER preserved, got %v", out)
	}
	if len(result.Added) != 2 {
		t.Errorf("expected 2 stripped, got %d", len(result.Added))
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "OTHER" {
		t.Errorf("expected OTHER skipped, got %v", result.Skipped)
	}
}

func TestStripPrefix_NilSource(t *testing.T) {
	_, _, err := StripPrefix(nil, "APP_")
	if err == nil {
		t.Error("expected error for nil source")
	}
}

func TestFormatPrefixResult_Output(t *testing.T) {
	r := PrefixResult{
		Added:   []string{"APP_FOO"},
		Skipped: []string{"APP_BAR"},
	}
	out := FormatPrefixResult(r)
	if out == "" {
		t.Error("expected non-empty output")
	}
}
