package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <env>",
		Short: "Export decrypted variables as shell export statements",
		Args:  cobra.ExactArgs(1),
		RunE:  runExport,
	}
	cmd.Flags().StringP("format", "f", "shell", "Output format: shell, dotenv, json")
	return cmd
}

func runExport(cmd *cobra.Command, args []string) error {
	env := args[0]
	format, _ := cmd.Flags().GetString("format")

	ks, err := keystore.New(keystoreDir())
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	st, err := store.New(storeDir())
	if err != nil {
		return fmt.Errorf("open store: %w", err)
	}

	data, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	env_data, err := envelope.Open(data, identity)
	if err != nil {
		return fmt.Errorf("decrypt env: %w", err)
	}

	vars, err := dotenv.Parse(strings.NewReader(string(env_data)))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case "shell":
		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "export %s=%q\n", k, vars[k])
		}
	case "dotenv":
		fmt.Fprint(os.Stdout, string(dotenv.Marshal(vars)))
	case "json":
		fmt.Fprint(os.Stdout, marshalJSON(keys, vars))
	default:
		return fmt.Errorf("unknown format %q: must be shell, dotenv, or json", format)
	}

	return nil
}

func marshalJSON(keys []string, vars map[string]string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		sb.WriteString(fmt.Sprintf("  %q: %q", k, vars[k]))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}
