package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tmc/envseal/internal/dotenv"
	"github.com/tmc/envseal/internal/envelope"
	"github.com/tmc/envseal/internal/keystore"
	"github.com/tmc/envseal/internal/store"
)

func newRenameKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renamekey <env> <old-key> <new-key>",
		Short: "Rename a key inside a sealed environment",
		Args:  cobra.ExactArgs(3),
		RunE:  runRenameKey,
	}
	return cmd
}

func runRenameKey(cmd *cobra.Command, args []string) error {
	env, oldKey, newKey := args[0], args[1], args[2]

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
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plaintext, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	parsed, err := dotenv.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("parse dotenv: %w", err)
	}

	renamed, result, err := dotenv.RenameKey(parsed, oldKey, newKey)
	if err != nil {
		return fmt.Errorf("rename key: %w", err)
	}

	marshaled := dotenv.Marshal(renamed)

	recipient := identity.Recipient()
	newSealed, err := envelope.Seal([]byte(marshaled), recipient)
	if err != nil {
		return fmt.Errorf("re-encrypt: %w", err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write sealed env: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), dotenv.FormatRenameResult(result))
	return nil
}
