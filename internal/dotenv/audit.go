package dotenv

import (
	"fmt"
	"sort"
	"time"
)

// AuditEventKind describes the type of audit event.
type AuditEventKind string

const (
	AuditSealed   AuditEventKind = "sealed"
	AuditOpened   AuditEventKind = "opened"
	AuditRekeyed  AuditEventKind = "rekeyed"
	AuditRotated  AuditEventKind = "rotated"
	AuditImported AuditEventKind = "imported"
	AuditExported AuditEventKind = "exported"
)

// AuditEntry records a single operation on an environment.
type AuditEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Kind      AuditEventKind `json:"kind"`
	Env       string         `json:"env"`
	Keys      []string       `json:"keys,omitempty"`
	Note      string         `json:"note,omitempty"`
}

// AuditLog is an ordered list of audit entries.
type AuditLog []AuditEntry

// NewEntry creates a new AuditEntry for the given environment and key map.
func NewEntry(kind AuditEventKind, env string, vars map[string]string, note string) AuditEntry {
	keys := Keys(vars)
	sort.Strings(keys)
	return AuditEntry{
		Timestamp: time.Now().UTC(),
		Kind:      kind,
		Env:       env,
		Keys:      keys,
		Note:      note,
	}
}

// FormatAuditLog returns a human-readable representation of the audit log.
func FormatAuditLog(log AuditLog) string {
	if len(log) == 0 {
		return "(no audit entries)\n"
	}
	out := ""
	for _, e := range log {
		line := fmt.Sprintf("%s  %-10s  env=%-12s  keys=%d",
			e.Timestamp.Format(time.RFC3339),
			e.Kind,
			e.Env,
			len(e.Keys),
		)
		if e.Note != "" {
			line += "  note=" + e.Note
		}
		out += line + "\n"
	}
	return out
}
