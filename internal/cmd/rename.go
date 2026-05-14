package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRenameCmd(keystoreDir, storeDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <old-env> <new-env>",
		Short: "Rename a sealed environment",
		Long: `Rename a sealed environment and its associated key.

This copies the sealed file and key to the new name, then removes the originals.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRename(cmd, args, keystoreDir(), storeDir())
		},
	}
	return cmd
}

func runRename(cmd *cobra.Command, args []string, ksDir, stDir string) error {
	oldEnv := args[0]
	newEnv := args[1]

	if oldEnv == newEnv {
		return fmt.Errorf("old and new environment names must differ")
	}

	ks, err := newKeystore(ksDir)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	st, err := newStore(stDir)
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	// Check source exists
	if !ks.Exists(oldEnv) {
		return fmt.Errorf("no key found for environment %q", oldEnv)
	}

	data, err := st.Read(oldEnv)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", oldEnv, err)
	}

	identity, err := ks.Load(oldEnv)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", oldEnv, err)
	}

	// Write to new destination
	if err := st.Write(newEnv, data); err != nil {
		return fmt.Errorf("write sealed env %q: %w", newEnv, err)
	}

	if err := ks.Save(newEnv, identity); err != nil {
		return fmt.Errorf("save key for %q: %w", newEnv, err)
	}

	// Remove originals
	if err := st.Delete(oldEnv); err != nil {
		return fmt.Errorf("delete old sealed env %q: %w", oldEnv, err)
	}

	if err := ks.Delete(oldEnv); err != nil {
		return fmt.Errorf("delete old key %q: %w", oldEnv, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "renamed environment %q to %q\n", oldEnv, newEnv)
	return nil
}
