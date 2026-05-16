package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envseal/internal/dotenv"
	"envseal/internal/keystore"
	"envseal/internal/store"
)

func newGrepCmd() *cobra.Command {
	var (
		keysOnly   bool
		valuesOnly bool
		ignoreCase bool
		invert     bool
	)

	cmd := &cobra.Command{
		Use:   "grep <pattern> <env>",
		Short: "Search keys and values in a sealed environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGrep(cmd, args, keysOnly, valuesOnly, ignoreCase, invert)
		},
	}

	cmd.Flags().BoolVar(&keysOnly, "keys", false, "Search only keys")
	cmd.Flags().BoolVar(&valuesOnly, "values", false, "Search only values")
	cmd.Flags().BoolVarP(&ignoreCase, "ignore-case", "i", false, "Case-insensitive matching")
	cmd.Flags().BoolVarP(&invert, "invert", "v", false, "Invert match (show non-matching entries)")
	return cmd
}

func runGrep(cmd *cobra.Command, args []string, keysOnly, valuesOnly, ignoreCase, invert bool) error {
	pattern, env := args[0], args[1]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}

	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	envMap, err := openEnvelope(sealed, identity)
	if err != nil {
		return fmt.Errorf("open envelope: %w", err)
	}

	opts := dotenv.GrepOptions{
		Pattern:      pattern,
		SearchKeys:   keysOnly,
		SearchValues: valuesOnly,
		IgnoreCase:   ignoreCase,
		Invert:       invert,
	}

	results, err := dotenv.Grep(envMap, opts)
	if err != nil {
		return fmt.Errorf("grep: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "(no matches)")
		return nil
	}

	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "%s=%s\n", r.Key, r.Value)
	}
	fmt.Fprint(cmd.OutOrStdout(), sb.String())
	return nil
}
