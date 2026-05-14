package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yourorg/envseal/internal/keystore"
)

func TestInitCmd_NoArgs(t *testing.T) {
	cmd := newInitCmd(t.TempDir(), t.TempDir())
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
}

func TestInitCmd_CreatesKey(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	cmd := newInitCmd(ksDir, stDir)
	cmd.SetArgs([]string{"development"})
	err := cmd.Execute()
	require.NoError(t, err)

	ks, err := keystore.New(ksDir)
	require.NoError(t, err)

	exists, err := ks.Exists("development")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestInitCmd_Idempotent(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()

	for i := 0; i < 2; i++ {
		cmd := newInitCmd(ksDir, stDir)
		cmd.SetArgs([]string{"staging"})
		err := cmd.Execute()
		require.NoError(t, err)
	}

	ks, err := keystore.New(ksDir)
	require.NoError(t, err)

	exists, err := ks.Exists("staging")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestInitCmd_MultipleEnvs(t *testing.T) {
	ksDir := t.TempDir()
	stDir := t.TempDir()
	envs := []string{"development", "staging", "production"}

	for _, env := range envs {
		cmd := newInitCmd(ksDir, stDir)
		cmd.SetArgs([]string{env})
		err := cmd.Execute()
		require.NoError(t, err)
	}

	ks, err := keystore.New(ksDir)
	require.NoError(t, err)

	for _, env := range envs {
		exists, err := ks.Exists(env)
		require.NoError(t, err)
		require.True(t, exists, "key for %s should exist", env)
	}
}
