package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/keystore"
)

func newKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage environment keys",
	}
	cmd.AddCommand(newKeysListCmd())
	cmd.AddCommand(newKeysDeleteCmd())
	return cmd
}

func newKeysListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List environments with stored keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := keystore.New(keystoreDir())
			if err != nil {
				return fmt.Errorf("open keystore: %w", err)
			}
			envs, err := ks.List()
			if err != nil {
				return fmt.Errorf("list keys: %w", err)
			}
			if len(envs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(no keys stored)")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(envs, "\n"))
			return nil
		},
	}
}

func newKeysDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <environment>",
		Short: "Delete the key for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env := args[0]
			ks, err := keystore.New(keystoreDir())
			if err != nil {
				return fmt.Errorf("open keystore: %w", err)
			}
			if err := ks.Delete(env); err != nil {
				return fmt.Errorf("delete key for %q: %w", env, err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "deleted key for environment %q\n", env)
			return nil
		},
	}
}
