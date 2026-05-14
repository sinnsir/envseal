package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newRotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate <environment>",
		Short: "Re-encrypt a sealed env file with a new key",
		Args:  cobra.ExactArgs(1),
		RunE:  runRotate,
	}
	return cmd
}

func runRotate(cmd *cobra.Command, args []string) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir())
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	oldIdentity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	st, err := store.New(storeDir())
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plain, err := envelope.Open(sealed, oldIdentity)
	if err != nil {
		return fmt.Errorf("decrypt %q: %w", env, err)
	}

	newIdentity, err := ks.Generate(env)
	if err != nil {
		return fmt.Errorf("generate new key: %w", err)
	}

	newSealed, err := envelope.Seal(plain, newIdentity.Recipient())
	if err != nil {
		return fmt.Errorf("re-encrypt %q: %w", env, err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write sealed env: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "rotated key for environment %q\n", env)
	return nil
}
