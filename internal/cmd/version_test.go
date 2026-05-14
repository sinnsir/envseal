package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersionCmd_Output(t *testing.T) {
	cmd := newVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.NoError(t, err)
	out := buf.String()
	require.True(t, strings.Contains(out, "envseal"), "output should contain app name")
}

func TestVersionCmd_ContainsVersion(t *testing.T) {
	cmd := newVersionCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	err := cmd.Execute()
	require.NoError(t, err)
	out := buf.String()
	// Should contain some version-like string
	require.NotEmpty(t, strings.TrimSpace(out))
}
