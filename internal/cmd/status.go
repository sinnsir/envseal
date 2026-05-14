package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envseal/internal/dotenv"
	"github.com/nicholasgasior/envseal/internal/envelope"
	"github.com/nicholasgasior/envseal/internal/keystore"
	"github.com/nicholasgasior/envseal/internal/store"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status <env>",
		Short: "Show status of a sealed environment vs a plain .env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runStatus,
	}
	cmd.Flags().StringP("file", "f", ".env", "Path to the plain .env file")
	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	env := args[0]
	plainFile, _ := cmd.Flags().GetString("file")

	ks, err := keystore.New(keystoreDir())
	if err != nil {
		return fmt.Errorf("opening keystore: %w", err)
	}
	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("loading key for %q: %w", env, err)
	}

	s, err := store.New(storeDir())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	sealedData, err := s.Read(env)
	if err != nil {
		return fmt.Errorf("reading sealed env %q: %w", env, err)
	}

	env_map, err := envelope.Open(sealedData, identity)
	if err != nil {
		return fmt.Errorf("decrypting: %w", err)
	}
	sealedVars, err := dotenv.Parse(string(env_map))
	if err != nil {
		return fmt.Errorf("parsing sealed env: %w", err)
	}

	plainBytes, err := os.ReadFile(plainFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", plainFile, err)
	}
	plainVars, err := dotenv.Parse(string(plainBytes))
	if err != nil {
		return fmt.Errorf("parsing plain env: %w", err)
	}

	diffs := dotenv.Diff(sealedVars, plainVars)
	if len(diffs) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No changes between plain file and sealed environment.")
		return nil
	}

	fmt.Fprintln(cmd.OutOrStdout(), dotenv.FormatDiff(diffs))
	return nil
}
