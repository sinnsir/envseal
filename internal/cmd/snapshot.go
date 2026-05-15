package cmd

import (
	"fmt"

	"github.com/nicholasgasior/envseal/internal/dotenv"
	"github.com/nicholasgasior/envseal/internal/envelope"
	"github.com/nicholasgasior/envseal/internal/keystore"
	"github.com/nicholasgasior/envseal/internal/store"
	"github.com/spf13/cobra"
)

func newSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot <environment>",
		Short: "Print a stable fingerprint of a sealed environment",
		Long: `Decrypts the sealed environment and prints a SHA-256 snapshot
of its key/value pairs. Useful for auditing or detecting drift
without exposing the actual values.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runSnapshot,
		Example: "  envseal snapshot production",
	}
	return cmd
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed file for %q: %w", env, err)
	}

	env_bytes, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt %q: %w", env, err)
	}

	m, err := dotenv.Parse(string(env_bytes))
	if err != nil {
		return fmt.Errorf("parse dotenv: %w", err)
	}

	snap := dotenv.TakeSnapshot(m)
	fmt.Fprintln(cmd.OutOrStdout(), snap.String())
	return nil
}
