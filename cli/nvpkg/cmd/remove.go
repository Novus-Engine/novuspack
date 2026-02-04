package cmd

import (
	"context"
	"fmt"
	"os"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var removePattern bool

var removeCmd = &cobra.Command{
	Use:   "remove <package path> <internal path or pattern>",
	Short: "Remove a file or directory from a NovusPack package",
	Long:  "Removes a single file, all files under a directory path (path ending with /), or files matching a pattern (--pattern).",
	Args:  cobra.ExactArgs(2),
	RunE:  runRemove,
}

func init() {
	removeCmd.Flags().BoolVar(&removePattern, "pattern", false, "treat second argument as a glob pattern (e.g. *.tmp)")
}

func runRemove(_ *cobra.Command, args []string) error {
	pkgPath := args[0]
	internalPath := args[1]
	ctx := context.Background()

	pkg, err := novuspack.OpenPackage(ctx, pkgPath)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	switch {
	case removePattern:
		_, err := pkg.RemoveFilePattern(ctx, internalPath)
		if err != nil {
			return fmt.Errorf("remove pattern: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Removed files matching %q from %s\n", internalPath, pkgPath)
	case internalPath != "" && internalPath[len(internalPath)-1] == '/':
		_, err := pkg.RemoveDirectory(ctx, internalPath, nil)
		if err != nil {
			return fmt.Errorf("remove directory: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Removed directory %s from %s\n", internalPath, pkgPath)
	default:
		if err := pkg.RemoveFile(ctx, internalPath); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Removed %s from %s\n", internalPath, pkgPath)
	}
	if err := pkg.Write(ctx); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
