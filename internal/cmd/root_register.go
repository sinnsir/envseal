package cmd

func registerSubcommands(root *rootCmd) {
	root.cmd.AddCommand(newInitCmd(root.keystoreDir))
	root.cmd.AddCommand(newSealCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newOpenCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newEditCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newDiffCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newRotateCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newCopyCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newRenameCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newExportCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newListCmd(root.storeDir))
	root.cmd.AddCommand(newStatusCmd(root.keystoreDir, root.storeDir))
	root.cmd.AddCommand(newKeysCmd(root.keystoreDir))
	root.cmd.AddCommand(newVersionCmd())
}
