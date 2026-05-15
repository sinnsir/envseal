package cmd

func registerSubcommands(root *rootCmd) {
	root.cmd.AddCommand(
		newInitCmd(),
		newSealCmd(),
		newOpenCmd(),
		newDiffCmd(),
		newRotateCmd(),
		newKeysCmd(),
		newVersionCmd(),
		newEditCmd(),
		newExportCmd(),
		newListCmd(),
		newStatusCmd(),
		newCopyCmd(),
		newRenameCmd(),
		newRekeyCmd(),
		newImportCmd(),
	)
}
