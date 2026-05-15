package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newEnvCmd() *cobra.Command {
	var keys []string

	cmd := &cobra.Command{
		Use:   "env <environment>",
		Short: "Print decrypted variables as shell exports",
		Long: `Decrypt the sealed environment and print each variable as an
export statement. Pipe to eval to apply them in the current shell:

  eval $(envseal env production)

Use --keys to restrict which variables are printed.`,
		Args:    cobra.ExactArgs(1),
		RunE:    func(cmd *cobra.Command, args []string) error { return runEnv(cmd, args, keys) },
		Example: "  envseal env production\n  envseal env staging --keys DB_URL,API_KEY",
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "comma-separated list of keys to include")
	return cmd
}

func runEnv(cmd *cobra.Command, args []string, keys []string) error {
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
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plaintext, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	m, err := dotenv.Parse(strings.NewReader(string(plaintext)))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	if len(keys) > 0 {
		m = dotenv.FromEnv(keys) // reuse filter logic via a local helper
		// override: filter the parsed map instead
		filtered := make(map[string]string, len(keys))
		for _, k := range keys {
			if v, ok := m[k]; ok {
				filtered[k] = v
			}
		}
		m = filtered
	}

	for _, pair := range dotenv.ToEnv(m) {
		fmt.Fprintf(cmd.OutOrStdout(), "export %s\n", pair)
	}
	return nil
}
