package dotenv

import (
	"testing"
)

func TestRedact_SensitiveKeys(t *testing.T) {
	input := map[string]string{
		"DATABASE_PASSWORD": "supersecret",
		"API_KEY":           "abc123",
		"SECRET_TOKEN":      "tok_xyz",
		"PRIVATE_KEY":       "-----BEGIN",
		"APP_NAME":          "myapp",
		"PORT":              "8080",
	}

	got := Redact(input)

	sensitive := []string{"DATABASE_PASSWORD", "API_KEY", "SECRET_TOKEN", "PRIVATE_KEY"}
	for _, k := range sensitive {
		if got[k] != RedactedValue {
			t.Errorf("key %q: expected %q, got %q", k, RedactedValue, got[k])
		}
	}

	plain := []string{"APP_NAME", "PORT"}
	for _, k := range plain {
		if got[k] != input[k] {
			t.Errorf("key %q: expected %q, got %q", k, input[k], got[k])
		}
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{
		"PASSWORD": "original",
		"HOST":     "localhost",
	}

	_ = Redact(input)

	if input["PASSWORD"] != "original" {
		t.Errorf("Redact mutated the input map")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	got := Redact(map[string]string{})
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRedact_CaseInsensitive(t *testing.T) {
	input := map[string]string{
		"db_password": "secret",
		"Auth_Token":  "tok",
		"normal_var":  "value",
	}

	got := Redact(input)

	if got["db_password"] != RedactedValue {
		t.Errorf("expected db_password to be redacted")
	}
	if got["Auth_Token"] != RedactedValue {
		t.Errorf("expected Auth_Token to be redacted")
	}
	if got["normal_var"] != "value" {
		t.Errorf("expected normal_var to be unchanged")
	}
}
