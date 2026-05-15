package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newMergeCmd() *cobra.Command {
	var strategy string

	cmd := &cobra.Command{
		Use:   "merge <src-env> <dst-env>",
		Short: "Merge secrets from one environment into another",
		Long: `Merge decrypts both environments and merges the source into the destination.

Strategies:
  overwrite     - src keys replace dst keys on conflict (default)
  keep-existing - dst keys are preserved on conflict
  error         - abort if any key conflicts`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var strat dotenv.MergeStrategy
			switch strategy {
			case "overwrite":
				strat = dotenv.StrategyOverwrite
			case "keep-existing":
				strat = dotenv.StrategyKeepExisting
			case "error":
				strat = dotenv.StrategyError
			default:
				return fmt.Errorf("unknown strategy %q; use overwrite, keep-existing, or error", strategy)
			}
			return runMerge(cmd, args[0], args[1], strat)
		},
	}

	cmd.Flags().StringVar(&strategy, "strategy", "overwrite", "conflict resolution strategy (overwrite|keep-existing|error)")
	return cmd
}

func runMerge(cmd *cobra.Command, srcEnv, dstEnv string, strategy dotenv.MergeStrategy) error {
	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}
	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	decrypt := func(env string) (map[string]string, error) {
		identity, err := ks.Load(env)
		if err != nil {
			return nil, fmt.Errorf("load key for %q: %w", env, err)
		}
		sealed, err := st.Read(env)
		if err != nil {
			return nil, fmt.Errorf("read sealed env %q: %w", env, err)
		}
		env_map, err := envelope.Open(sealed, identity)
		if err != nil {
			return nil, fmt.Errorf("decrypt %q: %w", env, err)
		}
		return dotenv.Parse(string(env_map))
	}

	srcMap, err := decrypt(srcEnv)
	if err != nil {
		return err
	}
	dstMap, err := decrypt(dstEnv)
	if err != nil {
		return err
	}

	merged, err := dotenv.MergeWithStrategy(dstMap, srcMap, strategy)
	if err != nil {
		return fmt.Errorf("merge conflict: %w", err)
	}

	dstIdentity, err := ks.Load(dstEnv)
	if err != nil {
		return fmt.Errorf("load dst key: %w", err)
	}
	plain := dotenv.Marshal(merged)
	sealed, err := envelope.Seal([]byte(plain), dstIdentity.Recipient())
	if err != nil {
		return fmt.Errorf("seal merged env: %w", err)
	}
	if err := st.Write(dstEnv, sealed); err != nil {
		return fmt.Errorf("write merged env: %w", err)
	}

	fmt.Fprintf(os.Stdout, "merged %q into %q (%d keys)\n", srcEnv, dstEnv, len(merged))
	return nil
}
