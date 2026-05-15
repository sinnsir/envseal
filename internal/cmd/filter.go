package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newFilterCmd() *cobra.Command {
	var prefix, suffix, pattern string
	var invert bool

	cmd := &cobra.Command{
		Use:   "filter <env>",
		Short: "Print a filtered subset of keys from a sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFilter(cmd, args[0], dotenv.FilterOptions{
				Prefix:  prefix,
				Suffix:  suffix,
				Pattern: pattern,
				Invert:  invert,
			})
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "retain keys with this prefix")
	cmd.Flags().StringVar(&suffix, "suffix", "", "retain keys with this suffix")
	cmd.Flags().StringVarP(&pattern, "pattern", "p", "", "retain keys matching this regex")
	cmd.Flags().BoolVarP(&invert, "invert", "v", false, "invert the filter (exclude matching keys)")
	return cmd
}

func runFilter(cmd *cobra.Command, env string, opts dotenv.FilterOptions) error {
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

	plain, err := sealed.Open(identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	filtered, err := dotenv.Filter(plain, opts)
	if err != nil {
		return fmt.Errorf("filter: %w", err)
	}

	out := dotenv.Marshal(filtered)
	_, err = fmt.Fprint(cmd.OutOrStdout(), out)
	return err
}
