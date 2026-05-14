package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEditCmd_NoArgs(t *testing.T) {
	root := newRootCmd()
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"edit"})
	err := root.Execute()
	require.Error(t, err)
}

func TestEditCmd_MissingSealedEnv(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE_DIR", filepath.Join(dir, "keys"))
	t.Setenv("ENVSEAL_STORE_DIR", filepath.Join(dir, "store"))

	// init a key for "staging"
	root := newRootCmd()
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"init", "staging"})
	require.NoError(t, root.Execute())

	// edit without a sealed file present should error
	root2 := newRootCmd()
	root2.SetOut(&bytes.Buffer{})
	root2.SetErr(&bytes.Buffer{})
	root2.SetArgs([]string{"edit", "staging"})
	err := root2.Execute()
	require.Error(t, err)
}

func TestEditCmd_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVSEAL_KEYSTORE_DIR", filepath.Join(dir, "keys"))
	t.Setenv("ENVSEAL_STORE_DIR", filepath.Join(dir, "store"))

	// write a fake editor script that appends a new variable
	editorScript := filepath.Join(dir, "editor.sh")
	err := os.WriteFile(editorScript, []byte("#!/bin/sh\necho 'EDITED=true' >> \"$1\"\n"), 0755)
	require.NoError(t, err)
	t.Setenv("EDITOR", editorScript)

	// init key
	root := newRootCmd()
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"init", "prod"})
	require.NoError(t, root.Execute())

	// seal initial file
	envFile := filepath.Join(dir, ".env.prod")
	err = os.WriteFile(envFile, []byte("FOO=bar\n"), 0600)
	require.NoError(t, err)

	root2 := newRootCmd()
	root2.SetOut(&bytes.Buffer{})
	root2.SetErr(&bytes.Buffer{})
	root2.SetArgs([]string{"seal", "prod", envFile})
	require.NoError(t, root2.Execute())

	// edit
	root3 := newRootCmd()
	out := &bytes.Buffer{}
	root3.SetOut(out)
	root3.SetErr(&bytes.Buffer{})
	root3.SetArgs([]string{"edit", "prod"})
	require.NoError(t, root3.Execute())
	require.Contains(t, out.String(), "sealed")

	// open and verify the edit was persisted
	outFile := filepath.Join(dir, ".env.prod.out")
	root4 := newRootCmd()
	root4.SetOut(&bytes.Buffer{})
	root4.SetErr(&bytes.Buffer{})
	root4.SetArgs([]string{"open", "prod", "--output", outFile})
	require.NoError(t, root4.Execute())

	contents, err := os.ReadFile(outFile)
	require.NoError(t, err)
	require.Contains(t, string(contents), "EDITED=true")
	require.Contains(t, string(contents), "FOO=bar")
}
