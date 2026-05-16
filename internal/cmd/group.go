package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tmc/envseal/internal/dotenv"
	"github.com/tmc/envseal/internal/envelope"
	"github.com/tmc/envseal/internal/keystore"
	"github.com/tmc/envseal/internal/store"
)

func newGroupCmd() *cobra.Command {
	var sep string
	var format string

	cmd := &cobra.Command{
		Use:   "group <env>",
		Short: "Group keys in a sealed environment by prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroup(cmd, args, sep, format)
		},
	}

	cmd.Flags().StringVar(&sep, "sep", "_", "Key separator for grouping")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text or list")
	return cmd
}

func runGroup(cmd *cobra.Command, args []string, sep, format string) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("group: open keystore: %w", err)
	}

	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("group: open store: %w", err)
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("group: load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("group: read sealed env %q: %w", env, err)
	}

	env_map, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("group: decrypt: %w", err)
	}

	parsed, err := dotenv.Parse(strings.NewReader(string(env_map)))
	if err != nil {
		return fmt.Errorf("group: parse: %w", err)
	}

	result, err := dotenv.Group(parsed, sep)
	if err != nil {
		return fmt.Errorf("group: %w", err)
	}

	switch format {
	case "list":
		for prefix, keys := range result.Groups {
			for k := range keys {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", prefix, k)
			}
		}
		for k := range result.Ungrouped {
			fmt.Fprintf(cmd.OutOrStdout(), "ungrouped\t%s\n", k)
		}
	default:
		fmt.Fprint(cmd.OutOrStdout(), dotenv.FormatGroup(result))
	}

	return nil
}
