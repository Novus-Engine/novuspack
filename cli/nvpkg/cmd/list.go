package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <package path>",
	Short: "List files in a NovusPack package",
	Args:  cobra.ExactArgs(1),
	RunE:  runList,
}

var listReadOnly bool

func init() {
	listCmd.Flags().BoolVar(&listReadOnly, "read-only", false, "Open package read-only (no write risk)")
}

func runList(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	pkg, err := openPackage(ctx, path, listReadOnly)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	files, err := pkg.ListFiles()
	if err != nil {
		return fmt.Errorf("list files: %w", err)
	}
	for _, f := range files {
		_, _ = fmt.Fprintf(os.Stdout, "%s  %d  %d\n", f.PrimaryPath, f.Size, f.StoredSize)
	}
	return nil
}
