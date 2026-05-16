package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envseal/internal/dotenv"
	"envseal/internal/keystore"
	"envseal/internal/store"
)

func newTypecheckCmd() *cobra.Command {
	var hints []string

	cmd := &cobra.Command{
		Use:   "typecheck <env>",
		Short: "Validate value types in a sealed environment",
		Long: `Decrypt a sealed .env and validate values against expected types.

Type hints are provided as KEY=TYPE pairs, e.g.:
  envseal typecheck production --hint PORT=int --hint DEBUG=bool --hint API_URL=url

Supported types: string, int, float, bool, url, email`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTypecheck(cmd, args, hints)
		},
	}

	cmd.Flags().StringArrayVar(&hints, "hint", nil, "type hint as KEY=TYPE (repeatable)")
	_ = cmd.MarkFlagRequired("hint")
	return cmd
}

func runTypecheck(cmd *cobra.Command, args []string, rawHints []string) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}
	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	if !ks.Exists(env) {
		return fmt.Errorf("no key found for environment %q", env)
	}
	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key: %w", err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env: %w", err)
	}

	plain, err := sealed.Open(identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	parsed, err := dotenv.Parse(strings.NewReader(string(plain)))
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	hintMap := make(map[string]dotenv.TypeHint, len(rawHints))
	for _, h := range rawHints {
		parts := strings.SplitN(h, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid hint %q: expected KEY=TYPE", h)
		}
		hintMap[parts[0]] = dotenv.TypeHint(parts[1])
	}

	results := dotenv.TypeCheck(parsed, hintMap)
	fmt.Fprint(os.Stdout, dotenv.FormatTypeCheck(results))

	for _, r := range results {
		if !r.Valid {
			return fmt.Errorf("type check failed")
		}
	}
	return nil
}
