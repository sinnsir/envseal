package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

// BuildDate is set at build time via ldflags.
var BuildDate = "unknown"

// Commit is set at build time via ldflags.
var Commit = "none"

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the envseal version",
		Long:  `Print the current version, commit, and build date of envseal.`,
		RunE:  runVersion,
	}
	return cmd
}

func runVersion(cmd *cobra.Command, _ []string) error {
	fmt.Fprintf(cmd.OutOrStdout(), "envseal %s (commit: %s, built: %s)\n", Version, Commit, BuildDate)
	return nil
}
