package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/tmc/envseal/internal/dotenv"
	"github.com/tmc/envseal/internal/envelope"
	"github.com/tmc/envseal/internal/keystore"
	"github.com/tmc/envseal/internal/store"
)

func newCompareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare <env1> <env2>",
		Short: "Compare two sealed environments side by side",
		Args:  cobra.ExactArgs(2),
		RunE:  runCompare,
	}
	return cmd
}

func runCompare(cmd *cobra.Command, args []string) error {
	env1, env2 := args[0], args[1]

	ks, err := keystore.New(keystoreDir(cmd))
	if err != nil {
		return fmt.Errorf("keystore: %w", err)
	}
	st, err := store.New(storeDir(cmd))
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}

	open := func(env string) (map[string]string, error) {
		identity, err := ks.Load(env)
		if err != nil {
			return nil, fmt.Errorf("load key for %q: %w", env, err)
		}
		data, err := st.Read(env)
		if err != nil {
			return nil, fmt.Errorf("read sealed %q: %w", env, err)
		}
		env2data, err := envelope.Unmarshal(data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %q: %w", env, err)
		}
		plain, err := envelope.Open(env2data, identity)
		if err != nil {
			return nil, fmt.Errorf("open %q: %w", env, err)
		}
		return dotenv.Parse(string(plain))
	}

	map1, err := open(env1)
	if err != nil {
		return err
	}
	map2, err := open(env2)
	if err != nil {
		return err
	}

	result := dotenv.Compare(map1, map2)
	fmt.Fprintf(cmd.OutOrStdout(), "Comparing %s → %s\n", env1, env2)
	fmt.Fprintf(cmd.OutOrStdout(), "Summary: %s\n\n", result.Summary())

	printKeys := func(label string, keys map[string]string) {
		sorted := make([]string, 0, len(keys))
		for k := range keys {
			sorted = append(sorted, k)
		}
		sort.Strings(sorted)
		for _, k := range sorted {
			fmt.Fprintf(cmd.OutOrStdout(), "  %s %s\n", label, k)
		}
	}
	printKeys("+", result.Added)
	printKeys("-", result.Removed)

	changed := make([]string, 0, len(result.Changed))
	for k := range result.Changed {
		changed = append(changed, k)
	}
	sort.Strings(changed)
	for _, k := range changed {
		fmt.Fprintf(cmd.OutOrStdout(), "  ~ %s\n", k)
	}
	return nil
}
