package cmd

// registerSubcommands attaches all subcommands to the root command.
// This file centralises command registration so that new commands
// only need to be added in one place.
func registerSubcommands(root *rootCmd) {
	root.AddCommand(
		newInitCmd(),
		newSealCmd(),
		newOpenCmd(),
		newDiffCmd(),
		newEditCmd(),
		newExportCmd(),
		newRotateCmd(),
		newListCmd(),
		newStatusCmd(),
		newKeysCmd(),
		newVersionCmd(),
	)
}
