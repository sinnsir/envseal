package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envseal",
		Short: "Encrypt and version .env files using age encryption",
		Long: `envseal encrypts .env files using age encryption with per-environment
key management. Sealed files are git-friendly and can be safely committed.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.AddCommand(
		newInitCmd(),
		newSealCmd(),
		newOpenCmd(),
		newDiffCmd(),
		newRotateCmd(),
		newKeysCmd(),
		newVersionCmd(),
		newEditCmd(),
		newExportCmd(),
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
			ks, err := keystore.New(keystoreDir())
			if err != nil {
				return fmt.Errorf("open keystore: %w", err)
			}
			if ks.Exists(env) {
				fmt.Fprintf(cmd.OutOrStdout(), "key for %q already exists, skipping\n", env)
				return nil
			}
			identity, err := ks.Generate(env)
			if err != nil {
				return fmt.Errorf("generate key: %w", err)
			}
			recipient, err := identity.Recipient()
			if err != nil {
				return fmt.Errorf("derive recipient: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "generated key for %q\n", env)
			fmt.Fprintf(cmd.OutOrStdout(), "public key: %s\n", recipient)
			return nil
		},
	}
}

func keystoreDir() string {
	if v := os.Getenv("ENVSEAL_KEYSTORE"); v != "" {
		return v
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".envseal", "keys")
}

func storeDir() string {
	if v := os.Getenv("ENVSEAL_STORE"); v != "" {
		return v
	}
	return store.Default()
}
