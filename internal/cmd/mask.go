package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envseal/internal/dotenv"
	"envseal/internal/keystore"
	"envseal/internal/store"
)

func newMaskCmd(keystoreDir, storeDir func() string) *cobra.Command {
	var mode string
	var keys []string
	var exclude []string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "mask <env>",
		Short: "Display a masked view of a sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMask(cmd, args, keystoreDir(), storeDir(), mode, keys, exclude, verbose)
		},
	}
	cmd.Flags().StringVar(&mode, "mode", "full", "mask mode: full, partial, length")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "specific keys to mask (default: auto-detect sensitive)")
	cmd.Flags().StringSliceVar(&exclude, "exclude", nil, "keys to exclude from masking")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show which keys were masked")
	return cmd
}

func runMask(cmd *cobra.Command, args []string, keystoreDir, storeDir, modeStr string, keys, exclude []string, verbose bool) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}
	st, err := store.New(storeDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
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
		return fmt.Errorf("decrypt env %q: %w", env, err)
	}

	parsed, err := dotenv.Parse(strings.NewReader(string(plain)))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	var maskMode dotenv.MaskMode
	switch modeStr {
	case "partial":
		maskMode = dotenv.MaskPartial
	case "length":
		maskMode = dotenv.MaskLength
	default:
		maskMode = dotenv.MaskFull
	}

	masked := dotenv.Mask(parsed, dotenv.MaskOptions{
		Mode:    maskMode,
		Keys:    keys,
		Exclude: exclude,
	})

	if verbose {
		cmd.Print(dotenv.FormatMask(parsed, masked))
	}

	out, err := dotenv.Marshal(masked)
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}
	cmd.Print(string(out))
	return nil
}
