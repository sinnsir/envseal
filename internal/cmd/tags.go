package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envseal/internal/dotenv"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/store"
)

func newTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage tags on sealed environment variables",
	}
	cmd.AddCommand(newTagsListCmd())
	cmd.AddCommand(newTagsFilterCmd())
	return cmd
}

func newTagsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <env>",
		Short: "List tags attached to variables in a sealed environment",
		Args:  cobra.ExactArgs(1),
		RunE:  runTagsList,
	}
}

func runTagsList(cmd *cobra.Command, args []string) error {
	env := args[0]
	ks := keystore.New(keystoreDir())
	st := store.New(storeDir())

	id, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}
	data, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}
	envMap, err := envelope.Open(data, id)
	if err != nil {
		return fmt.Errorf("open envelope: %w", err)
	}

	// Tags are stored as a special key __TAGS__<VARNAME>=key=val,key=val
	for k, v := range envMap {
		if strings.HasPrefix(k, "__TAGS__") {
			varName := strings.TrimPrefix(k, "__TAGS__")
			tags, err := dotenv.ParseTags(v)
			if err != nil {
				continue
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", varName, dotenv.FormatTags(tags))
		}
	}
	return nil
}

func newTagsFilterCmd() *cobra.Command {
	var tagExpr string
	c := &cobra.Command{
		Use:   "filter <env>",
		Short: "Print variables matching a tag expression (key or key=value)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTagsFilter(cmd, args, tagExpr)
		},
	}
	c.Flags().StringVarP(&tagExpr, "tag", "t", "", "Tag expression: key or key=value (required)")
	_ = c.MarkFlagRequired("tag")
	return c
}

func runTagsFilter(cmd *cobra.Command, args []string, tagExpr string) error {
	env := args[0]
	ks := keystore.New(keystoreDir())
	st := store.New(storeDir())

	id, err := ks.Load(env)
	if err != nil {
		return fmt.Errorf("load key for %q: %w", env, err)
	}
	data, err := st.Read(env)
	if err != nil {
		return fmt.Errorf("read sealed env %q: %w", env, err)
	}
	envMap, err := envelope.Open(data, id)
	if err != nil {
		return fmt.Errorf("open envelope: %w", err)
	}

	tm := dotenv.TagMap{}
	plain := map[string]string{}
	for k, v := range envMap {
		if strings.HasPrefix(k, "__TAGS__") {
			varName := strings.TrimPrefix(k, "__TAGS__")
			tags, _ := dotenv.ParseTags(v)
			tm[varName] = tags
		} else {
			plain[k] = v
		}
	}

	tagKey, tagValue, _ := strings.Cut(tagExpr, "=")
	matched := dotenv.FilterByTag(plain, tm, tagKey, tagValue)
	for _, k := range dotenv.TaggedKeys(tm) {
		if v, ok := matched[k]; ok {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
	}
	return nil
}
