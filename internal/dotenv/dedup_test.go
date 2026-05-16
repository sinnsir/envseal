package dotenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDedup_NoDuplicates(t *testing.T) {
	pairs := []string{"FOO=bar", "BAZ=qux"}
	r, err := Dedup(pairs, DedupKeepFirst)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"FOO": "bar", "BAZ": "qux"}, r.Out)
	assert.Empty(t, r.Removed)
}

func TestDedup_KeepFirst(t *testing.T) {
	pairs := []string{"FOO=first", "BAR=one", "FOO=second"}
	r, err := Dedup(pairs, DedupKeepFirst)
	require.NoError(t, err)
	assert.Equal(t, "first", r.Out["FOO"])
	assert.Equal(t, []string{"FOO"}, r.Removed)
}

func TestDedup_KeepLast(t *testing.T) {
	pairs := []string{"FOO=first", "BAR=one", "FOO=second"}
	r, err := Dedup(pairs, DedupKeepLast)
	require.NoError(t, err)
	assert.Equal(t, "second", r.Out["FOO"])
	assert.Equal(t, []string{"FOO"}, r.Removed)
}

func TestDedup_MultipleDuplicates(t *testing.T) {
	pairs := []string{"A=1", "B=2", "A=3", "B=4", "C=5"}
	r, err := Dedup(pairs, DedupKeepFirst)
	require.NoError(t, err)
	assert.Equal(t, "1", r.Out["A"])
	assert.Equal(t, "2", r.Out["B"])
	assert.Equal(t, "5", r.Out["C"])
	assert.Equal(t, []string{"A", "B"}, r.Removed)
}

func TestDedup_SkipsComments(t *testing.T) {
	pairs := []string{"# comment", "FOO=bar"}
	r, err := Dedup(pairs, DedupKeepFirst)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"FOO": "bar"}, r.Out)
	assert.Empty(t, r.Removed)
}

func TestDedup_NilInput(t *testing.T) {
	_, err := Dedup(nil, DedupKeepFirst)
	require.Error(t, err)
}

func TestDedup_InvalidPair(t *testing.T) {
	_, err := Dedup([]string{"NOEQUALSSIGN"}, DedupKeepFirst)
	require.Error(t, err)
}

func TestFormatDedup_NoRemovals(t *testing.T) {
	r := &DedupResult{Out: map[string]string{"A": "1"}, Removed: nil}
	out := FormatDedup(r)
	assert.Contains(t, out, "no duplicate")
}

func TestFormatDedup_WithRemovals(t *testing.T) {
	r := &DedupResult{
		Out:     map[string]string{"A": "1"},
		Removed: []string{"A", "B"},
	}
	out := FormatDedup(r)
	assert.Contains(t, out, "2 duplicate")
	assert.Contains(t, out, "- A")
	assert.Contains(t, out, "- B")
}
