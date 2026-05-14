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

func newSealCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seal [env]",
		Short: "Encrypt a .env file for the given environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runSeal,
	}
	cmd.Flags().StringP("file", "f", "", "Path to the .env file to seal (default: .env.<env>)")
	return cmd
}

func runSeal(cmd *cobra.Command, args []string) error {
	env := args[0]

	filePath, _ := cmd.Flags().GetString("file")
	if filePath == "" {
		filePath = fmt.Sprintf(".env.%s", env)
	}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	parsed, err := dotenv.Parse(raw)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", filePath, err)
	}

	ks, err := keystore.New(keystoreDir())
	if err != nil {
		return fmt.Errorf("opening keystore: %w", err)
	}

	recipient, err := ks.Recipient(env)
	if err != nil {
		return fmt.Errorf("loading recipient for %q: %w", env, err)
	}

	plaintext := dotenv.Marshal(parsed)
	sealed, err := envelope.Seal(plaintext, recipient)
	if err != nil {
		return fmt.Errorf("sealing envelope: %w", err)
	}

	st, err := store.New(store.Default())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	if err := st.Write(env, sealed); err != nil {
		return fmt.Errorf("writing sealed file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "sealed %d keys into %s\n", len(parsed), env)
	return nil
}
