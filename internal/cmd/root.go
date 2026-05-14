package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const defaultKeystoreDir = ".envseal/keys"

var rootKeystoreDir string

// Execute builds and runs the root command.
func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envseal",
		Short: "Encrypt and version .env files using age encryption",
		Long: `envseal encrypts .env files per environment using age encryption.

Keys are stored locally in .envseal/keys and never committed to git.
Sealed files live in .envseal/ and are safe to commit.`,
		SilenceUsage: true,
	}

	root.PersistentFlags().StringVar(
		&rootKeystoreDir,
		"keystore",
		"",
		"Path to keystore directory (default: .envseal/keys)",
	)

	root.AddCommand(
		newSealCmd(),
		newOpenCmd(),
		newInitCmd(),
	)

	return root
}

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init [env]",
		Short: "Generate a new age key pair for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env := args[0]
			ks, err := keystore.New(keystoreDir())
			if err != nil {
				return err
			}
			pub, err := ks.Generate(env)
			if err != nil {
				return err
			}
			cmd.Printf("generated key for %q\npublic key: %s\n", env, pub)
			return nil
		},
	}
}

func keystoreDir() string {
	if rootKeystoreDir != "" {
		return rootKeystoreDir
	}
	if dir, err := os.Getwd(); err == nil {
		return filepath.Join(dir, defaultKeystoreDir)
	}
	return defaultKeystoreDir
}
