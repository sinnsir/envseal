package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/envseal/internal/dotenv"
	"github.com/yourusername/envseal/internal/envelope"
	"github.com/yourusername/envseal/internal/keystore"
	"github.com/yourusername/envseal/internal/store"
)

func newImportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <env> <file>",
		Short: "Import a plain .env file and seal it for the given environment",
		Args:  cobra.ExactArgs(2),
		RunE:  runImport,
	}
	cmd.Flags().Bool("overwrite", false, "Overwrite existing sealed env if present")
	return cmd
}

func runImport(cmd *cobra.Command, args []string) error {
	env := args[0]
	filePath := args[1]

	overwrite, _ := cmd.Flags().GetBool("overwrite")

	ks := keystore.New(keystoreDir())
	if !ks.Exists(env) {
		return fmt.Errorf("no key found for environment %q — run 'envseal init %s' first", env, env)
	}

	st := store.New(storeDir())
	if !overwrite {
		if _, err := st.Read(env); err == nil {
			return fmt.Errorf("sealed env %q already exists; use --overwrite to replace it", env)
		}
	}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	kvs, err := dotenv.Parse(string(raw))
	if err != nil {
		return fmt.Errorf("parsing .env file: %w", err)
	}

	recipient, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("loading key: %w", err)
	}

	sealed, err := envelope.Seal(kvs, recipient)
	if err != nil {
		return fmt.Errorf("sealing: %w", err)
	}

	data, err := envelope.Marshal(sealed)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	if err := st.Write(env, data); err != nil {
		return fmt.Errorf("writing sealed env: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "imported %d keys into sealed env %q\n", len(kvs), env)
	return nil
}
