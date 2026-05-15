package dotenv_test

import (
	"testing"

	"github.com/yourusername/envseal/internal/dotenv"
)

func TestPromote_AddsNewKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}
	out, res, err := dotenv.Promote(src, dst, dotenv.PromoteSkipExisting)
	if err != nil {
		t.Fatal(err)
	}
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected output: %v", out)
	}
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(res.Added))
	}
}

func TestPromote_SkipExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res, err := dotenv.Promote(src, dst, dotenv.PromoteSkipExisting)
	if err != nil {
		t.Fatal(err)
	}
	if out["A"] != "old" {
		t.Errorf("expected old value to be kept, got %q", out["A"])
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(res.Skipped))
	}
}

func TestPromote_Overwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, res, err := dotenv.Promote(src, dst, dotenv.PromoteOverwrite)
	if err != nil {
		t.Fatal(err)
	}
	if out["A"] != "new" {
		t.Errorf("expected new value, got %q", out["A"])
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
}

func TestPromote_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"B": "2"}
	_, _, _ = dotenv.Promote(src, dst, dotenv.PromoteOverwrite)
	if _, ok := dst["A"]; ok {
		t.Error("promote mutated dst")
	}
}

func TestPromote_NilSrcError(t *testing.T) {
	_, _, err := dotenv.Promote(nil, map[string]string{}, dotenv.PromoteSkipExisting)
	if err == nil {
		t.Error("expected error for nil src")
	}
}

func TestFormatPromoteResult_NoChanges(t *testing.T) {
	res := dotenv.PromoteResult{}
	got := dotenv.FormatPromoteResult(res)
	if got != "no changes\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestFormatPromoteResult_WithChanges(t *testing.T) {
	res := dotenv.PromoteResult{
		Added:       []string{"NEW_KEY"},
		Overwritten: []string{"CHANGED_KEY"},
		Skipped:     []string{"KEPT_KEY"},
	}
	got := dotenv.FormatPromoteResult(res)
	for _, want := range []string{"+ NEW_KEY", "~ CHANGED_KEY", "= KEPT_KEY (skipped)"} {
		if !contains(got, want) {
			t.Errorf("expected %q in output:\n%s", want, got)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
