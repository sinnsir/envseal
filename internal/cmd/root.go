package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

// Execute runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envseal",
		Short: "Encrypt and version .env files using age encryption",
	}
	root.AddCommand(
		newInitCmd(),
		newSealCmd(),
		newOpenCmd(),
		newDiffCmd(),
	)
	return root
}

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init <env>",
		Short: "Generate a new age key for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env := args[0]
			ksDir := keystoreDir()
			ks, err := keystore.New(ksDir)
			if err != nil {
				return fmt.Errorf("keystore: %w", err)
			}
			recipient, err := ks.Generate(env)
			if err != nil {
				return fmt.Errorf("generate key: %w", err)
			}
			fmt.Printf("Generated key for %q\nPublic key: %s\n", env, recipient)
			return nil
		},
	}
}

// keystoreDir returns the directory used for key storage,
// preferring the ENVSEAL_KEYSTORE env var, then ~/.config/envseal/keys.
func keystoreDir() string {
	if d := os.Getenv("ENVSEAL_KEYSTORE"); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".envseal/keys"
	}
	return filepath.Join(home, ".config", "envseal", "keys")
}

// storeDir returns the store directory, preferring ENVSEAL_STORE env var.
func storeDir() (string, error) {
	if d := os.Getenv("ENVSEAL_STORE"); d != "" {
		st, err := store.New(d)
		if err != nil {
			return "", err
		}
		_ = st
		return d, nil
	}
	return "", nil
}
