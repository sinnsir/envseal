package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envseal/internal/store"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all sealed environments",
		Long:  "List all sealed environments stored in the envseal store.",
		RunE:  runList,
	}
	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	dir := storeDir()
	s, err := store.New(dir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	envs, err := s.List()
	if err != nil {
		return fmt.Errorf("listing environments: %w", err)
	}

	if len(envs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No sealed environments found.")
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ENVIRONMENT")
	for _, env := range envs {
		fmt.Fprintln(w, env)
	}
	_ = w.Flush()
	_ = os.Stderr
	return nil
}
