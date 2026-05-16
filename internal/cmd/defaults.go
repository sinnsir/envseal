package cmd

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/envseal/internal/dotenv"
	"github.com/spf13/cobra"
)

func newDefaultsCmd() *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "defaults <env> <defaults-file>",
		Short: "Apply default values to a sealed environment",
		Long: `Reads a sealed environment and a plain .env defaults file.
Any key present in the defaults file but missing from the sealed env
is added. Existing keys are never overwritten.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDefaults(cmd, args, verbose)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show which keys were applied or skipped")
	return cmd
}

func runDefaults(cmd *cobra.Command, args []string, verbose bool) error {
	env := args[0]
	defaultsFile := args[1]

	ks, st, err := openKeystoreAndStore(cmd)
	if err != nil {
		return err
	}

	id, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	current, err := openEnvelope(sealed, id)
	if err != nil {
		return fmt.Errorf("decrypt env %q: %w", env, err)
	}

	defaultVals, err := dotenv.ReadFile(defaultsFile)
	if err != nil {
		return fmt.Errorf("read defaults file: %w", err)
	}

	merged, result, err := dotenv.ApplyDefaults(current, defaultVals)
	if err != nil {
		return fmt.Errorf("apply defaults: %w", err)
	}

	if err := sealAndWrite(merged, env, id, st); err != nil {
		return fmt.Errorf("seal env %q: %w", env, err)
	}

	if verbose {
		fmt.Fprint(os.Stdout, dotenv.FormatDefaults(result))
	} else {
		fmt.Fprintf(os.Stdout, "applied %d default(s) to %q\n", len(result.Applied), env)
	}
	return nil
}
