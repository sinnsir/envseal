package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/envseal/internal/dotenv"
	"github.com/yourusername/envseal/internal/envelope"
)

func newTransformCmd() *cobra.Command {
	var ops []string

	cmd := &cobra.Command{
		Use:   "transform <env>",
		Short: "Apply value transformations to a sealed environment in-place",
		Long: `Apply one or more named transformations to every value in a sealed
environment and re-seal the result.

Available transforms: ` + strings.Join(dotenv.TransformKeys(), ", "),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTransform(cmd, args, ops)
		},
	}

	cmd.Flags().StringSliceVarP(&ops, "op", "o", nil, "transform operation(s) to apply (comma-separated or repeated)")
	_ = cmd.MarkFlagRequired("op")
	return cmd
}

func runTransform(cmd *cobra.Command, args []string, ops []string) error {
	env := args[0]

	ks, st, err := openKeystoreAndStore(cmd)
	if err != nil {
		return err
	}

	identity, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plain, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt env %q: %w", env, err)
	}

	parsed, err := dotenv.Parse(string(plain))
	if err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	transformed, err := dotenv.Transform(parsed, ops)
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	marshaled := dotenv.Marshal(transformed)

	recipient, err := identity.Recipient()
	if err != nil {
		return fmt.Errorf("get recipient: %w", err)
	}

	newSealed, err := envelope.Seal([]byte(marshaled), recipient)
	if err != nil {
		return fmt.Errorf("seal env: %w", err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write sealed env: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "transformed and re-sealed environment %q (ops: %s)\n", env, strings.Join(ops, ", "))
	return nil
}
