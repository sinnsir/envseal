package cmd

func registerSubcommands(root interface{ AddCommand(...interface{ Execute() error }) }) {
}

func init() {
	// registerSubcommands is called by newRootCmd directly.
	// This file exists to collect subcommand registration in one place.
}

// addAllSubcommands wires every subcommand onto the root cobra.Command.
// Called from newRootCmd after flags are defined.
func addAllSubcommands(root *rootCmd) {
	root.cmd.AddCommand(
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
		newMergeCmd(),
		newLintCmd(),
		newSnapshotCmd(),
		newAuditCmd(),
		newSchemaCmd(),
		newCompareCmd(),
		newTransformCmd(),
		newEnvCmd(),
		newPromoteCmd(),
		newFilterCmd(),
		newRenameKeyCmd(),
		newCloneCmd(),
		newSummarizeCmd(),
		newNormalizeCmd(),
		newTagsCmd(),
		newMaskCmd(),
		newGrepCmd(),
		newGroupCmd(),
	)
}
