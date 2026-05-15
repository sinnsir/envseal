package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newSchemaCmd(ks *keystore.Keystore, st *store.Store) *cobra.Command {
	var schemaFile string
	cmd := &cobra.Command{
		Use:   "schema <env>",
		Short: "Validate a sealed environment against a schema file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchema(cmd, args, ks, st, schemaFile)
		},
	}
	cmd.Flags().StringVarP(&schemaFile, "schema", "s", ".env.schema", "path to schema definition file")
	return cmd
}

func runSchema(cmd *cobra.Command, args []string, ks *keystore.Keystore, st *store.Store, schemaFile string) error {
	env := args[0]

	// Load and decrypt the sealed environment.
	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}
	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}
	plaintext, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt env %q: %w", env, err)
	}
	envMap, err := dotenv.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	// Parse the schema file.
	rules, err := parseSchemaFile(schemaFile)
	if err != nil {
		return fmt.Errorf("parse schema file: %w", err)
	}

	result := dotenv.ValidateSchema(envMap, rules)
	fmt.Fprintln(cmd.OutOrStdout(), dotenv.FormatSchema(result))
	if result.HasIssues() {
		return fmt.Errorf("schema validation failed")
	}
	return nil
}

// parseSchemaFile reads a simple schema definition.
// Each line: KEY [required] [pattern=<substr>]
func parseSchemaFile(path string) ([]dotenv.SchemaRule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []dotenv.SchemaRule
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		rule := dotenv.SchemaRule{Key: parts[0]}
		for _, p := range parts[1:] {
			switch {
			case p == "required":
				rule.Required = true
			case strings.HasPrefix(p, "pattern="):
				rule.Pattern = strings.TrimPrefix(p, "pattern=")
			}
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
