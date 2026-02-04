package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <package path>",
	Short: "Show package metadata and summary",
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

var infoReadOnly bool

func init() {
	infoCmd.Flags().BoolVar(&infoReadOnly, "read-only", false, "Open package read-only (no write risk)")
}

func runInfo(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	pkg, err := openPackage(ctx, path, infoReadOnly)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	info, err := pkg.GetInfo()
	if err != nil {
		return fmt.Errorf("get info: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Path:       %s\n", path)
	_, _ = fmt.Fprintf(os.Stdout, "File count: %d\n", info.FileCount)
	_, _ = fmt.Fprintf(os.Stdout, "Uncompressed size: %d\n", info.FilesUncompressedSize)
	_, _ = fmt.Fprintf(os.Stdout, "Compressed size:   %d\n", info.FilesCompressedSize)
	if info.VendorID != 0 || info.AppID != 0 {
		_, _ = fmt.Fprintf(os.Stdout, "Vendor ID:  %d\n", info.VendorID)
		_, _ = fmt.Fprintf(os.Stdout, "App ID:     %d\n", info.AppID)
	}
	if info.Comment != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Comment:    %s\n", info.Comment)
	}
	return nil
}
