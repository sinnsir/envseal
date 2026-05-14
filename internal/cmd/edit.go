package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit <env>",
		Short: "Decrypt, open in $EDITOR, then re-seal a .env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runEdit,
	}
	return cmd
}

func runEdit(cmd *cobra.Command, args []string) error {
	env := args[0]

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

	sealed, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}

	plaintext, err := envelope.Open(sealed, identity)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "envseal-edit-*.env")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(plaintext); err != nil {
		tmpFile.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	c := exec.Command(editor, tmpFile.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("editor exited with error: %w", err)
	}

	edited, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("read edited file: %w", err)
	}

	if _, err := dotenv.Parse(edited); err != nil {
		return fmt.Errorf("invalid .env syntax: %w", err)
	}

	newSealed, err := envelope.Seal(edited, identity.Recipient())
	if err != nil {
		return fmt.Errorf("re-seal: %w", err)
	}

	if err := st.Write(env, newSealed); err != nil {
		return fmt.Errorf("write sealed env: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "sealed %q successfully\n", env)
	return nil
}
