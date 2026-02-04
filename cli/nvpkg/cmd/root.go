package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nvpkg",
	Short: "NovusPack package manager CLI",
	Long:  "nvpkg is a CLI for creating, inspecting, and modifying NovusPack (.nvpk) packages. Use 'nvpkg interactive' (or 'nvpkg i') for REPL mode.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(headerCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(commentCmd)
	rootCmd.AddCommand(identityCmd)
	rootCmd.AddCommand(metadataCmd)
	rootCmd.AddCommand(interactiveCmd)
}
