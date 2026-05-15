package dotenv

import (
	"strings"
	"testing"
	"time"
)

func TestNewEntry_Fields(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	e := NewEntry(AuditSealed, "production", vars, "initial seal")

	if e.Kind != AuditSealed {
		t.Errorf("expected kind %q, got %q", AuditSealed, e.Kind)
	}
	if e.Env != "production" {
		t.Errorf("expected env %q, got %q", "production", e.Env)
	}
	if len(e.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(e.Keys))
	}
	if e.Note != "initial seal" {
		t.Errorf("unexpected note: %q", e.Note)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestNewEntry_KeysSorted(t *testing.T) {
	vars := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	e := NewEntry(AuditOpened, "staging", vars, "")

	for i := 1; i < len(e.Keys); i++ {
		if e.Keys[i-1] > e.Keys[i] {
			t.Errorf("keys not sorted: %v", e.Keys)
		}
	}
}

func TestFormatAuditLog_Empty(t *testing.T) {
	out := FormatAuditLog(AuditLog{})
	if !strings.Contains(out, "no audit entries") {
		t.Errorf("expected empty message, got: %q", out)
	}
}

func TestFormatAuditLog_WithEntries(t *testing.T) {
	log := AuditLog{
		{
			Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			Kind:      AuditSealed,
			Env:       "production",
			Keys:      []string{"DB_URL", "SECRET"},
			Note:      "deploy",
		},
		{
			Timestamp: time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC),
			Kind:      AuditRekeyed,
			Env:       "production",
			Keys:      []string{"DB_URL", "SECRET"},
		},
	}
	out := FormatAuditLog(log)
	if !strings.Contains(out, "sealed") {
		t.Error("expected 'sealed' in output")
	}
	if !strings.Contains(out, "rekeyed") {
		t.Error("expected 'rekeyed' in output")
	}
	if !strings.Contains(out, "production") {
		t.Error("expected 'production' in output")
	}
	if !strings.Contains(out, "note=deploy") {
		t.Error("expected note in output")
	}
}

func TestFormatAuditLog_NoNote(t *testing.T) {
	log := AuditLog{
		{Timestamp: time.Now(), Kind: AuditExported, Env: "dev", Keys: []string{"K"}},
	}
	out := FormatAuditLog(log)
	if strings.Contains(out, "note=") {
		t.Error("should not contain 'note=' when note is empty")
	}
}
