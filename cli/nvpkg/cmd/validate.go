package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <package path>",
	Short: "Validate package integrity",
	Long:  "Validates an existing NovusPack package (header, index, and optional content checks).",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

var validateReadOnly bool

func init() {
	validateCmd.Flags().BoolVar(&validateReadOnly, "read-only", false, "Open package read-only (no write risk)")
}

func runValidate(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	pkg, err := openPackage(ctx, path, validateReadOnly)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	if err := pkg.Validate(ctx); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "OK %s\n", path)
	return nil
}
