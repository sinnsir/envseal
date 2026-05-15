package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newRekeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rekey <env>",
		Short: "Re-encrypt a sealed env file with a newly generated key",
		Long: `Generates a new age key for the given environment, decrypts the
existing sealed file with the old key, re-encrypts it with the new key,
and replaces both the sealed file and the stored key.`,
		Args:    cobra.ExactArgs(1),
		RunE:    runRekey,
		Example: "  envseal rekey production",
	}
	return cmd
}

func runRekey(cmd *cobra.Command, args []string) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	// Load old identity to decrypt existing sealed file.
	oldIdentity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	// Read and decrypt the existing sealed file.
	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed file for %q: %w", env, err)
	}

	env_map, err := envelope.Open(sealed, oldIdentity)
	if err != nil {
		return fmt.Errorf("decrypt sealed file: %w", err)
	}

	// Generate a fresh key.
	newIdentity, err := ks.Generate(env)
	if err != nil {
		return fmt.Errorf("generate new key: %w", err)
	}

	// Re-encrypt with the new key.
	newSealed, err := envelope.Seal(env_map, newIdentity.Recipient())
	if err != nil {
		return fmt.Errorf("re-encrypt: %w", err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write sealed file: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "rekeyed environment %q\n", env)
	return nil
}
