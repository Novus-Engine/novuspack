package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <package path> <internal path>",
	Short: "Read a file from a NovusPack package to stdout or a file",
	Args:  cobra.ExactArgs(2),
	RunE:  runRead,
}

var readOutput string
var readReadOnly bool

func init() {
	readCmd.Flags().StringVarP(&readOutput, "output", "o", "", "Write to file instead of stdout")
	readCmd.Flags().BoolVar(&readReadOnly, "read-only", false, "Open package read-only (no write risk)")
}

func runRead(_ *cobra.Command, args []string) error {
	pkgPath := args[0]
	internalPath := args[1]
	ctx := context.Background()

	pkg, err := openPackage(ctx, pkgPath, readReadOnly)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	data, err := pkg.ReadFile(ctx, internalPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if readOutput != "" {
		if err := os.WriteFile(readOutput, data, 0o644); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
		return nil
	}
	_, err = os.Stdout.Write(data)
	return err
}
