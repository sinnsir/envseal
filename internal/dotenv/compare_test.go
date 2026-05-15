package dotenv_test

import (
	"testing"

	"github.com/tmc/envseal/internal/dotenv"
)

func TestCompare_Added(t *testing.T) {
	old := map[string]string{"A": "1"}
	new := map[string]string{"A": "1", "B": "2"}
	r := dotenv.Compare(old, new)
	if len(r.Added) != 1 || r.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got %v", r.Added)
	}
	if r.HasChanges() == false {
		t.Error("expected HasChanges to be true")
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1"}
	r := dotenv.Compare(old, new)
	if len(r.Removed) != 1 || r.Removed["B"] != "2" {
		t.Errorf("expected B=2 in Removed, got %v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	old := map[string]string{"A": "old"}
	new := map[string]string{"A": "new"}
	r := dotenv.Compare(old, new)
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(r.Changed))
	}
	pair := r.Changed["A"]
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected [old new], got %v", pair)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	r := dotenv.Compare(m, m)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
	if r.Summary() != "no changes" {
		t.Errorf("unexpected summary: %s", r.Summary())
	}
}

func TestCompare_Summary(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "updated", "C": "3"}
	r := dotenv.Compare(old, new)
	s := r.Summary()
	if s == "" || s == "no changes" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestCompare_Same(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1", "B": "changed"}
	r := dotenv.Compare(old, new)
	if len(r.Same) != 1 || r.Same["A"] != "1" {
		t.Errorf("expected A=1 in Same, got %v", r.Same)
	}
}
