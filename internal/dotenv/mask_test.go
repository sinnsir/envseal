package dotenv

import (
	"strings"
	"testing"
)

func TestMask_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"API_KEY":   "supersecret",
		"HOST":      "localhost",
		"DB_PASSWORD": "hunter2",
	}
	masked := Mask(env, MaskOptions{Mode: MaskFull})
	if masked["API_KEY"] != "***" {
		t.Errorf("expected API_KEY to be masked, got %q", masked["API_KEY"])
	}
	if masked["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD to be masked, got %q", masked["DB_PASSWORD"])
	}
	if masked["HOST"] != "localhost" {
		t.Errorf("expected HOST to be unmasked, got %q", masked["HOST"])
	}
}

func TestMask_ExplicitKeys(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	masked := Mask(env, MaskOptions{Mode: MaskFull, Keys: []string{"FOO"}})
	if masked["FOO"] != "***" {
		t.Errorf("expected FOO masked, got %q", masked["FOO"])
	}
	if masked["BAZ"] != "qux" {
		t.Errorf("expected BAZ unmasked, got %q", masked["BAZ"])
	}
}

func TestMask_ExcludeKey(t *testing.T) {
	env := map[string]string{"SECRET": "abc123"}
	masked := Mask(env, MaskOptions{Mode: MaskFull, Exclude: []string{"SECRET"}})
	if masked["SECRET"] != "abc123" {
		t.Errorf("expected SECRET excluded from masking, got %q", masked["SECRET"])
	}
}

func TestMask_PartialMode(t *testing.T) {
	env := map[string]string{"API_KEY": "abcdef"}
	masked := Mask(env, MaskOptions{Mode: MaskPartial})
	v := masked["API_KEY"]
	if !strings.HasPrefix(v, "a") || !strings.HasSuffix(v, "f") {
		t.Errorf("partial mask should preserve first/last char, got %q", v)
	}
	if !strings.Contains(v, "***") {
		t.Errorf("partial mask should contain asterisks, got %q", v)
	}
}

func TestMask_LengthMode(t *testing.T) {
	env := map[string]string{"API_SECRET": "hello"}
	masked := Mask(env, MaskOptions{Mode: MaskLength})
	if masked["API_SECRET"] != "*****" {
		t.Errorf("expected 5 asterisks, got %q", masked["API_SECRET"])
	}
}

func TestMask_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"API_KEY": "original"}
	_ = Mask(env, MaskOptions{Mode: MaskFull})
	if env["API_KEY"] != "original" {
		t.Error("Mask mutated the input map")
	}
}

func TestMask_EmptyValue(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	masked := Mask(env, MaskOptions{Mode: MaskFull})
	if masked["API_KEY"] != "" {
		t.Errorf("empty value should remain empty, got %q", masked["API_KEY"])
	}
}

func TestFormatMask_ContainsMaskedCount(t *testing.T) {
	original := map[string]string{"API_KEY": "secret", "HOST": "localhost"}
	masked := Mask(original, MaskOptions{Mode: MaskFull})
	out := FormatMask(original, masked)
	if !strings.Contains(out, "masked 1 key(s)") {
		t.Errorf("expected masked count in output, got: %s", out)
	}
}
