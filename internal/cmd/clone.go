package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yourusername/envseal/internal/dotenv"
	"github.com/yourusername/envseal/internal/envelope"
	"github.com/yourusername/envseal/internal/keystore"
	"github.com/yourusername/envseal/internal/store"
)

func newCloneCmd() *cobra.Command {
	var keys []string

	cmd := &cobra.Command{
		Use:   "clone <src-env> <dst-env>",
		Short: "Clone a sealed environment into a new environment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClone(cmd, args, keys)
		},
	}

	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "restrict clone to specific keys (comma-separated)")
	return cmd
}

func runClone(cmd *cobra.Command, args []string, keys []string) error {
	srcEnv, dstEnv := args[0], args[1]
	if srcEnv == dstEnv {
		return fmt.Errorf("source and destination environments must differ")
	}

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}

	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	// Load source identity
	srcIdentity, err := ks.Load(srcEnv)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", srcEnv, err)
	}

	// Read and decrypt source sealed file
	raw, err := st.Read(srcEnv)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", srcEnv, err)
	}

	env, err := envelope.Open(raw, srcIdentity)
	if err != nil {
		return fmt.Errorf("open envelope: %w", err)
	}

	parsed, err := dotenv.Parse(string(env))
	if err != nil {
		return fmt.Errorf("parse dotenv: %w", err)
	}

	// Clone into subset if keys specified
	cloned, result, err := dotenv.Clone(parsed, keys)
	if err != nil {
		return fmt.Errorf("clone: %w", err)
	}

	// Ensure destination key exists
	if !ks.Exists(dstEnv) {
		if _, err := ks.Generate(dstEnv); err != nil {
			return fmt.Errorf("generate key for %q: %w", dstEnv, err)
		}
	}

	dstIdentity, err := ks.Load(dstEnv)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", dstEnv, err)
	}

	sealed, err := envelope.Seal([]byte(dotenv.Marshal(cloned)), dstIdentity.Recipient())
	if err != nil {
		return fmt.Errorf("seal: %w", err)
	}

	if err := st.Write(dstEnv, sealed); err != nil {
		return fmt.Errorf("write sealed env: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), dotenv.FormatCloneResult(result))
	return nil
}
