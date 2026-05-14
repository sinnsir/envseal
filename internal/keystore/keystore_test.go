package keystore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndLoad(t *testing.T) {
	ks, err := New(t.TempDir())
	require.NoError(t, err)

	id, err := ks.Generate("staging")
	require.NoError(t, err)
	require.NotNil(t, id)

	loaded, err := ks.Load("staging")
	require.NoError(t, err)
	require.Equal(t, id.String(), loaded.String())
}

func TestLoad_RoundTrip(t *testing.T) {
	ks, err := New(t.TempDir())
	require.NoError(t, err)

	_, err = ks.Generate("prod")
	require.NoError(t, err)

	id2, err := ks.Generate("prod") // overwrite
	require.NoError(t, err)

	loaded, err := ks.Load("prod")
	require.NoError(t, err)
	require.Equal(t, id2.String(), loaded.String())
}

func TestLoad_NotFound(t *testing.T) {
	ks, err := New(t.TempDir())
	require.NoError(t, err)

	_, err = ks.Load("missing")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestExists(t *testing.T) {
	ks, err := New(t.TempDir())
	require.NoError(t, err)

	require.False(t, ks.Exists("dev"))
	_, err = ks.Generate("dev")
	require.NoError(t, err)
	require.True(t, ks.Exists("dev"))
}

func TestDelete(t *testing.T) {
	ks, err := New(t.TempDir())
	require.NoError(t, err)

	_, err = ks.Generate("dev")
	require.NoError(t, err)

	require.NoError(t, ks.Delete("dev"))
	require.False(t, ks.Exists("dev"))

	err = ks.Delete("dev")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestList(t *testing.T) {
	ks, err := New(t.TempDir())
	require.NoError(t, err)

	envs := []string{"dev", "staging", "prod"}
	for _, e := range envs {
		_, err = ks.Generate(e)
		require.NoError(t, err)
	}

	list, err := ks.List()
	require.NoError(t, err)
	require.ElementsMatch(t, envs, list)
}
