package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/envseal/internal/dotenv"
	"github.com/yourusername/envseal/internal/envelope"
	"github.com/yourusername/envseal/internal/keystore"
	"github.com/yourusername/envseal/internal/store"
)

func newSummarizeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "summarize <env>",
		Short: "Print a summary of keys in a sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runSummarize,
	}
}

func runSummarize(cmd *cobra.Command, args []string) error {
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

	parsed, err := dotenv.Parse(string(plaintext))
	if err != nil {
		return fmt.Errorf("parse dotenv: %w", err)
	}

	summary := dotenv.Summarize(parsed)
	cmd.Print(dotenv.FormatSummary(summary))
	return nil
}
