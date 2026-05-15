package dotenv_test

import (
	"testing"

	"github.com/nicholasgasior/envseal/internal/dotenv"
)

func TestTakeSnapshot_Deterministic(t *testing.T) {
	m := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	s1 := dotenv.TakeSnapshot(m)
	s2 := dotenv.TakeSnapshot(m)
	if s1.Hash != s2.Hash {
		t.Errorf("expected identical hashes, got %q and %q", s1.Hash, s2.Hash)
	}
}

func TestTakeSnapshot_OrderIndependent(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAR": "2", "FOO": "1"}
	sa := dotenv.TakeSnapshot(a)
	sb := dotenv.TakeSnapshot(b)
	if sa.Hash != sb.Hash {
		t.Errorf("expected same hash regardless of map iteration order")
	}
}

func TestTakeSnapshot_DifferentValues(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}
	sa := dotenv.TakeSnapshot(a)
	sb := dotenv.TakeSnapshot(b)
	if sa.Equal(sb) {
		t.Error("expected different snapshots for different values")
	}
}

func TestTakeSnapshot_EmptyMap(t *testing.T) {
	s := dotenv.TakeSnapshot(map[string]string{})
	if s.Count != 0 {
		t.Errorf("expected count 0, got %d", s.Count)
	}
	if len(s.Hash) != 64 {
		t.Errorf("expected 64-char hex hash, got len %d", len(s.Hash))
	}
}

func TestSnapshot_Short(t *testing.T) {
	s := dotenv.TakeSnapshot(map[string]string{"A": "b"})
	if len(s.Short()) != 12 {
		t.Errorf("expected Short() to return 12 chars, got %d", len(s.Short()))
	}
}

func TestSnapshot_Equal(t *testing.T) {
	m := map[string]string{"X": "y"}
	s1 := dotenv.TakeSnapshot(m)
	s2 := dotenv.TakeSnapshot(m)
	if !s1.Equal(s2) {
		t.Error("expected Equal to return true for identical maps")
	}
}

func TestSnapshot_String(t *testing.T) {
	m := map[string]string{"FOO": "bar"}
	s := dotenv.TakeSnapshot(m)
	str := s.String()
	if str == "" {
		t.Error("expected non-empty String()")
	}
}

func TestTakeSnapshot_KeysAreSorted(t *testing.T) {
	m := map[string]string{"ZZZ": "1", "AAA": "2", "MMM": "3"}
	s := dotenv.TakeSnapshot(m)
	if len(s.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(s.Keys))
	}
	if s.Keys[0] != "AAA" || s.Keys[1] != "MMM" || s.Keys[2] != "ZZZ" {
		t.Errorf("keys not sorted: %v", s.Keys)
	}
}
