package dotenv

import (
	"strings"
	"testing"
)

func TestSummarize_BasicCounts(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"DB_PASSWORD": "secret",
		"EMPTY_VAL": "",
		"API_KEY": "abc123",
	}
	s := Summarize(env)
	if s.Total != 4 {
		t.Errorf("Total: got %d, want 4", s.Total)
	}
	if s.Empty != 1 {
		t.Errorf("Empty: got %d, want 1", s.Empty)
	}
	if s.Sensitive != 2 {
		t.Errorf("Sensitive: got %d, want 2 (DB_PASSWORD, API_KEY)", s.Sensitive)
	}
}

func TestSummarize_KeysSorted(t *testing.T) {
	env := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	s := Summarize(env)
	if len(s.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(s.Keys))
	}
	if s.Keys[0] != "A_KEY" || s.Keys[1] != "M_KEY" || s.Keys[2] != "Z_KEY" {
		t.Errorf("keys not sorted: %v", s.Keys)
	}
}

func TestSummarize_EmptyMap(t *testing.T) {
	s := Summarize(map[string]string{})
	if s.Total != 0 || s.Empty != 0 || s.Sensitive != 0 {
		t.Errorf("expected all zeros for empty map, got %+v", s)
	}
}

func TestFormatSummary_ContainsFields(t *testing.T) {
	s := Summary{Total: 3, Empty: 1, Sensitive: 1, Keys: []string{"A", "B", "C"}}
	out := FormatSummary(s)
	for _, want := range []string{"Total keys", "Empty values", "Sensitive", "- A", "- B", "- C"} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatSummary output missing %q\nGot:\n%s", want, out)
		}
	}
}

func TestFormatSummary_NoKeysSection(t *testing.T) {
	s := Summary{Total: 0, Empty: 0, Sensitive: 0, Keys: []string{}}
	out := FormatSummary(s)
	if strings.Contains(out, "Keys:") {
		t.Errorf("expected no Keys section for empty summary, got:\n%s", out)
	}
}
