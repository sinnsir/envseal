package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nicholasgasior/envseal/internal/dotenv"
	"github.com/nicholasgasior/envseal/internal/envelope"
	"github.com/nicholasgasior/envseal/internal/keystore"
	"github.com/nicholasgasior/envseal/internal/store"
	"github.com/spf13/cobra"
)

func newAuditCmd(ks *keystore.Keystore, st *store.Store) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "audit <env>",
		Short: "Show audit log of operations for an environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudit(cmd, ks, st, args[0], format)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	return cmd
}

func runAudit(cmd *cobra.Command, ks *keystore.Keystore, st *store.Store, env, format string) error {
	if !st.Exists(env) {
		return fmt.Errorf("no sealed file found for environment %q", env)
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed file: %w", err)
	}

	env2, vars, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	_ = env2

	entry := dotenv.NewEntry(dotenv.AuditOpened, env, vars, "audit query")
	log := dotenv.AuditLog{entry}

	switch format {
	case "json":
		return json.NewEncoder(os.Stdout).Encode(log)
	default:
		fmt.Fprint(cmd.OutOrStdout(), dotenv.FormatAuditLog(log))
	}
	return nil
}
