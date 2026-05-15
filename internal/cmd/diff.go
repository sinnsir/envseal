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

func newDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <env> <file>",
		Short: "Show diff between a sealed env and a plaintext .env file",
		Args:  cobra.ExactArgs(2),
		RunE:  runDiff,
	}
	return cmd
}

func runDiff(cmd *cobra.Command, args []string) error {
	env := args[0]
	plainPath := args[1]

	ksDir := keystoreDir()
	ks, err := keystore.New(ksDir)
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	st, err := store.Default()
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	env1Raw, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("open sealed env: %w", err)
	}

	oldVars, err := dotenv.Parse(string(env1Raw))
	if err != nil {
		return fmt.Errorf("parse sealed env: %w", err)
	}

	plainRaw, err := os.ReadFile(plainPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %q does not exist", plainPath)
		}
		return fmt.Errorf("read %q: %w", plainPath, err)
	}

	newVars, err := dotenv.Parse(string(plainRaw))
	if err != nil {
		return fmt.Errorf("parse %q: %w", plainPath, err)
	}

	changes := dotenv.Diff(oldVars, newVars)
	output := dotenv.FormatDiff(changes)
	if output == "" {
		fmt.Println("No changes.")
		return nil
	}
	fmt.Print(output)
	return nil
}
