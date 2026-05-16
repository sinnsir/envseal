package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nicholasgasior/envseal/internal/dotenv"
	"github.com/nicholasgasior/envseal/internal/envelope"
	"github.com/nicholasgasior/envseal/internal/keystore"
	"github.com/nicholasgasior/envseal/internal/store"
)

func newNormalizeCmd() *cobra.Command {
	var opts struct {
		uppercase    bool
		trimValues   bool
		removeEmpty  bool
		quoteValues  bool
		dryRun       bool
	}

	cmd := &cobra.Command{
		Use:   "normalize <env>",
		Short: "Normalize keys and values in a sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNormalize(cmd, args, opts.uppercase, opts.trimValues, opts.removeEmpty, opts.quoteValues, opts.dryRun)
		},
	}

	cmd.Flags().BoolVar(&opts.uppercase, "uppercase", false, "uppercase all keys")
	cmd.Flags().BoolVar(&opts.trimValues, "trim", false, "trim whitespace from values")
	cmd.Flags().BoolVar(&opts.removeEmpty, "remove-empty", false, "remove keys with empty values")
	cmd.Flags().BoolVar(&opts.quoteValues, "quote", false, "quote values containing whitespace")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "print changes without writing")

	return cmd
}

func runNormalize(cmd *cobra.Command, args []string, uppercase, trimValues, removeEmpty, quoteValues, dryRun bool) error {
	env := args[0]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}

	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	id, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	env_map, err := envelope.Open(sealed, id)
	if err != nil {
		return fmt.Errorf("open envelope: %w", err)
	}

	parsed, err := dotenv.Parse(string(env_map))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	var normalizeOpts []dotenv.NormalizeOption
	if uppercase {
		normalizeOpts = append(normalizeOpts, dotenv.NormalizeUppercaseKeys)
	}
	if trimValues {
		normalizeOpts = append(normalizeOpts, dotenv.NormalizeTrimValues)
	}
	if removeEmpty {
		normalizeOpts = append(normalizeOpts, dotenv.NormalizeRemoveEmpty)
	}
	if quoteValues {
		normalizeOpts = append(normalizeOpts, dotenv.NormalizeQuoteValues)
	}

	result, err := dotenv.Normalize(parsed, normalizeOpts)
	if err != nil {
		return fmt.Errorf("normalize: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), dotenv.FormatNormalizeResult(result))

	if dryRun {
		return nil
	}

	marshaled, err := dotenv.Marshal(result.Output)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	recipient, err := id.Recipient()
	if err != nil {
		return fmt.Errorf("recipient: %w", err)
	}

	newSealed, err := envelope.Seal([]byte(marshaled), recipient)
	if err != nil {
		return fmt.Errorf("seal: %w", err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}
