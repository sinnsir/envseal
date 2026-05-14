package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCopyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "copy <src-env> <dst-env>",
		Short: "Copy a sealed environment to a new environment",
		Long: `Copy decrypts a sealed environment using the source environment's key,
then re-encrypts it using the destination environment's key (generating one if needed).`,
		Args:    cobra.ExactArgs(2),
		RunE:    runCopy,
		Example: "  envseal copy production staging",
	}
	return cmd
}

func runCopy(cmd *cobra.Command, args []string) error {
	srcEnv := args[0]
	dstEnv := args[1]

	if srcEnv == dstEnv {
		return fmt.Errorf("source and destination environments must differ")
	}

	ks, err := keystoreDir(cmd)
	if err != nil {
		return err
	}

	st, err := storeDir(cmd)
	if err != nil {
		return err
	}

	// Load source identity
	srcIdentity, err := ks.Load(srcEnv)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", srcEnv, err)
	}

	// Read and decrypt source sealed file
	srcData, err := st.Read(srcEnv)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", srcEnv, err)
	}

	env, err := envelope.Open(srcData, srcIdentity)
	if err != nil {
		return fmt.Errorf("open sealed env %q: %w", srcEnv, err)
	}

	// Ensure destination key exists, generate if not
	var dstRecipient age.Recipient
	if ks.Exists(dstEnv) {
		dstIdentity, err := ks.Load(dstEnv)
		if err != nil {
			return fmt.Errorf("load key for %q: %w", dstEnv, err)
		}
		dstRecipient = dstIdentity.Recipient()
	} else {
		dstIdentity, err := ks.Generate(dstEnv)
		if err != nil {
			return fmt.Errorf("generate key for %q: %w", dstEnv, err)
		}
		dstRecipient = dstIdentity.Recipient()
		fmt.Fprintf(cmd.OutOrStdout(), "Generated new key for environment %q\n", dstEnv)
	}

	// Re-encrypt for destination
	dstData, err := envelope.Seal(env, dstRecipient)
	if err != nil {
		return fmt.Errorf("seal for %q: %w", dstEnv, err)
	}

	if err := st.Write(dstEnv, dstData); err != nil {
		return fmt.Errorf("write sealed env %q: %w", dstEnv, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Copied environment %q → %q\n", srcEnv, dstEnv)
	return nil
}
