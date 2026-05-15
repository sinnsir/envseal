package cmd

func registerSubcommands(root *rootCmd) {
	root.AddCommand(newInitCmd())
	root.AddCommand(newSealCmd())
	root.AddCommand(newOpenCmd())
	root.AddCommand(newDiffCmd())
	root.AddCommand(newEditCmd())
	root.AddCommand(newExportCmd())
	root.AddCommand(newListCmd())
	root.AddCommand(newStatusCmd())
	root.AddCommand(newRotateCmd())
	root.AddCommand(newCopyCmd())
	root.AddCommand(newRenameCmd())
	root.AddCommand(newRekeyCmd())
	root.AddCommand(newKeysCmd())
	root.AddCommand(newVersionCmd())
}
