package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata <package path>",
	Short: "Show full package metadata",
	Long:  "Prints package metadata (PackageInfo, file entries, path metadata). Use --json for machine-readable output.",
	Args:  cobra.ExactArgs(1),
	RunE:  runMetadata,
}

var metadataJSON bool
var metadataReadOnly bool

func init() {
	metadataCmd.Flags().BoolVar(&metadataJSON, "json", false, "Output as JSON")
	metadataCmd.Flags().BoolVar(&metadataReadOnly, "read-only", false, "Open package read-only")
}

func runMetadata(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	pkg, err := openPackage(ctx, path, metadataReadOnly)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	meta, err := pkg.GetMetadata()
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}
	if meta == nil {
		return fmt.Errorf("metadata is nil")
	}

	if metadataJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(meta); err != nil {
			return fmt.Errorf("json encode: %w", err)
		}
		return nil
	}
	_, _ = fmt.Fprintf(os.Stdout, "File entries: %d\n", len(meta.FileEntries))
	_, _ = fmt.Fprintf(os.Stdout, "Path metadata entries: %d\n", len(meta.PathMetadataEntries))
	if meta.PackageInfo != nil {
		_, _ = fmt.Fprintf(os.Stdout, "File count: %d\n", meta.FileCount)
		_, _ = fmt.Fprintf(os.Stdout, "Uncompressed size: %d\n", meta.FilesUncompressedSize)
		_, _ = fmt.Fprintf(os.Stdout, "Compressed size: %d\n", meta.FilesCompressedSize)
	}
	return nil
}
