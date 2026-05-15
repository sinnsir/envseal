package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yourusername/envseal/internal/dotenv"
	"github.com/yourusername/envseal/internal/keystore"
	"github.com/yourusername/envseal/internal/store"
)

func newPromoteCmd() *cobra.Command {
	var overwrite bool
	cmd := &cobra.Command{
		Use:   "promote <src-env> <dst-env>",
		Short: "Promote variables from one sealed environment into another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPromote(cmd, args, overwrite)
		},
	}
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in destination")
	return cmd
}

func runPromote(cmd *cobra.Command, args []string, overwrite bool) error {
	srcEnv, dstEnv := args[0], args[1]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}
	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	srcID, err := ks.Load(srcEnv)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", srcEnv, err)
	}
	dstID, err := ks.Load(dstEnv)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", dstEnv, err)
	}

	srcMap, err := openEnv(st, srcEnv, srcID)
	if err != nil {
		return fmt.Errorf("open %q: %w", srcEnv, err)
	}
	dstMap, err := openEnv(st, dstEnv, dstID)
	if err != nil {
		return fmt.Errorf("open %q: %w", dstEnv, err)
	}

	strategy := dotenv.PromoteSkipExisting
	if overwrite {
		strategy = dotenv.PromoteOverwrite
	}

	merged, result, err := dotenv.Promote(srcMap, dstMap, strategy)
	if err != nil {
		return err
	}

	if err := sealEnv(st, dstEnv, dstID, merged); err != nil {
		return fmt.Errorf("seal %q: %w", dstEnv, err)
	}

	fmt.Fprint(cmd.OutOrStdout(), dotenv.FormatPromoteResult(result))
	return nil
}
