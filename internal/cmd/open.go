package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newOpenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open [env]",
		Short: "Decrypt a sealed .env file for the given environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runOpen,
	}
	cmd.Flags().StringP("output", "o", "", "Output file path (default: .env.<env>)")
	cmd.Flags().BoolP("export", "e", false, "Print as export statements instead of writing a file")
	return cmd
}

func runOpen(cmd *cobra.Command, args []string) error {
	env := args[0]
	exportMode, _ := cmd.Flags().GetBool("export")
	outPath, _ := cmd.Flags().GetString("output")

	st, err := store.New(store.Default())
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("reading sealed file for %q: %w", env, err)
	}

	ks, err := keystore.New(keystoreDir())
	if err != nil {
		return fmt.Errorf("opening keystore: %w", err)
	}

	identity, err := ks.Identity(env)
	if err != nil {
		return fmt.Errorf("loading identity for %q: %w", env, err)
	}

	plaintext, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("opening envelope: %w", err)
	}

	if exportMode {
		parsed, err := dotenv.Parse(plaintext)
		if err != nil {
			return fmt.Errorf("parsing decrypted content: %w", err)
		}
		for k, v := range parsed {
			fmt.Fprintf(cmd.OutOrStdout(), "export %s=%q\n", k, v)
		}
		return nil
	}

	if outPath == "" {
		outPath = fmt.Sprintf(".env.%s", env)
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o700); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, plaintext, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", outPath, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "opened %s -> %s\n", env, outPath)
	return nil
}
