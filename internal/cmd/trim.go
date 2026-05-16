package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/envseal/internal/dotenv"
	"github.com/yourusername/envseal/internal/envelope"
	"github.com/yourusername/envseal/internal/keystore"
	"github.com/yourusername/envseal/internal/store"
)

func newTrimCmd(ks *keystore.KeyStore, st *store.Store) *cobra.Command {
	var (
		leading  bool
		trailing bool
		quotes   bool
		prefix   string
		suffix   string
		dryRun   bool
	)
	cmd := &cobra.Command{
		Use:   "trim <env>",
		Short: "Trim whitespace or quotes from env values",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTrim(cmd, args, ks, st, dotenv.TrimOptions{
				LeadingWhitespace:  leading,
				TrailingWhitespace: trailing,
				Quotes:             quotes,
				Prefix:             prefix,
				Suffix:             suffix,
			}, dryRun)
		},
	}
	cmd.Flags().BoolVar(&leading, "leading", false, "trim leading whitespace")
	cmd.Flags().BoolVar(&trailing, "trailing", false, "trim trailing whitespace")
	cmd.Flags().BoolVar(&quotes, "quotes", false, "strip surrounding quotes")
	cmd.Flags().StringVar(&prefix, "prefix", "", "strip value prefix")
	cmd.Flags().StringVar(&suffix, "suffix", "", "strip value suffix")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show changes without saving")
	return cmd
}

func runTrim(cmd *cobra.Command, args []string, ks *keystore.KeyStore, st *store.Store, opts dotenv.TrimOptions, dryRun bool) error {
	env := args[0]

	id, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	raw, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	env2, err := envelope.Open(raw, id)
	if err != nil {
		return fmt.Errorf("open envelope: %w", err)
	}

	m, err := dotenv.Parse(string(env2))
	if err != nil {
		return fmt.Errorf("parse dotenv: %w", err)
	}

	result, err := dotenv.Trim(m, opts)
	if err != nil {
		return fmt.Errorf("trim: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), dotenv.FormatTrim(result))

	if dryRun {
		return nil
	}

	sealed, err := envelope.Seal([]byte(dotenv.Marshal(result.Out)), id.Recipient())
	if err != nil {
		return fmt.Errorf("seal: %w", err)
	}
	return st.Write(args[0], sealed)
}
