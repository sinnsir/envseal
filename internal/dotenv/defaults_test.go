package dotenv_test

import (
	"testing"

	"github.com/nicholasgasior/envseal/internal/dotenv"
)

func TestApplyDefaults_FillsMissing(t *testing.T) {
	dst := map[string]string{"A": "1"}
	defaults := map[string]string{"A": "99", "B": "2", "C": "3"}

	out, result, err := dotenv.ApplyDefaults(dst, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %s", out["C"])
	}
	if len(result.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(result.Applied))
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "A" {
		t.Errorf("expected A skipped, got %v", result.Skipped)
	}
}

func TestApplyDefaults_DoesNotMutateDst(t *testing.T) {
	dst := map[string]string{"X": "original"}
	defaults := map[string]string{"Y": "new"}

	out, _, err := dotenv.ApplyDefaults(dst, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["Y"]; ok {
		t.Error("dst was mutated")
	}
	if out["Y"] != "new" {
		t.Errorf("expected Y=new in output")
	}
}

func TestApplyDefaults_NilDst(t *testing.T) {
	_, _, err := dotenv.ApplyDefaults(nil, map[string]string{"A": "1"})
	if err == nil {
		t.Error("expected error for nil dst")
	}
}

func TestApplyDefaults_NilDefaults(t *testing.T) {
	_, _, err := dotenv.ApplyDefaults(map[string]string{}, nil)
	if err == nil {
		t.Error("expected error for nil defaults")
	}
}

func TestApplyDefaults_EmptyDefaults(t *testing.T) {
	dst := map[string]string{"A": "1"}
	out, result, err := dotenv.ApplyDefaults(dst, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Applied) != 0 {
		t.Errorf("expected 0 applied")
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1")
	}
}

func TestFormatDefaults_Applied(t *testing.T) {
	r := dotenv.DefaultsResult{
		Applied:  map[string]string{"FOO": "bar"},
		Skipped: []string{"BAZ"},
	}
	s := dotenv.FormatDefaults(r)
	if s == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormatDefaults_Empty(t *testing.T) {
	r := dotenv.DefaultsResult{}
	s := dotenv.FormatDefaults(r)
	if s != "no defaults to apply\n" {
		t.Errorf("unexpected output: %q", s)
	}
}
