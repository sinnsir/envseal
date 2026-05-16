package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newTruncateCmd(ks *keystore.Keystore, st *store.Store) *cobra.Command {
	var maxLen int
	var suffix string
	var keys []string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "truncate <env>",
		Short: "Truncate long values in a sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTruncate(cmd, args, ks, st, maxLen, suffix, keys, dryRun)
		},
	}

	cmd.Flags().IntVar(&maxLen, "max-len", 80, "maximum value length before truncation")
	cmd.Flags().StringVar(&suffix, "suffix", "...", "suffix appended to truncated values")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "restrict truncation to these keys (comma-separated)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print changes without saving")

	return cmd
}

func runTruncate(
	cmd *cobra.Command,
	args []string,
	ks *keystore.Keystore,
	st *store.Store,
	maxLen int,
	suffix string,
	keys []string,
	dryRun bool,
) error {
	env := args[0]

	id, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plain, err := envelope.Open(sealed, id)
	if err != nil {
		return fmt.Errorf("decrypt %q: %w", env, err)
	}

	parsed, err := dotenv.Parse(plain)
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	res, err := dotenv.Truncate(parsed, dotenv.TruncateOptions{
		MaxLen: maxLen,
		Suffix: suffix,
		Keys:   keys,
	})
	if err != nil {
		return fmt.Errorf("truncate: %w", err)
	}

	cmd.Println(dotenv.FormatTruncate(res))

	if dryRun {
		return nil
	}

	updated, err := dotenv.Marshal(res.Output)
	if err != nil {
		return fmt.Errorf("marshal env: %w", err)
	}

	newSealed, err := envelope.Seal(updated, id.Recipient())
	if err != nil {
		return fmt.Errorf("re-encrypt %q: %w", env, err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write sealed env %q: %w", env, err)
	}
	return nil
}
