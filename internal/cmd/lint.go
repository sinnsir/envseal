package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/envseal/internal/dotenv"
	"github.com/your-org/envseal/internal/envelope"
	"github.com/your-org/envseal/internal/keystore"
	"github.com/your-org/envseal/internal/store"
)

func newLintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint <env>",
		Short: "Check a sealed environment for common style issues",
		Args:  cobra.ExactArgs(1),
		RunE:  runLint,
	}
	return cmd
}

func runLint(cmd *cobra.Command, args []string) error {
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

	data, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plaintext, err := envelope.Open(data, identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	parsed, err := dotenv.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	issues := dotenv.Lint(parsed)
	if len(issues) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "no issues found")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%d issue(s) found in %q:\n", len(issues), env)
	for _, issue := range issues {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", issue)
	}

	// Exit with a non-zero status so CI pipelines can catch lint failures.
	os.Exit(1)
	return nil
}
