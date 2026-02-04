package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var extractCmd = &cobra.Command{
	Use:   "extract <package path> [internal path]",
	Short: "Extract all or a subtree of files from a NovusPack package to a directory",
	Long:  "Extract files from a .nvpk package to the filesystem. With no internal path, extracts all files. With an internal path (e.g. /docs), extracts only that file or directory subtree. Destination is set with -o/--output.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runExtract,
}

var extractOutput string
var extractReadOnly bool

func init() {
	extractCmd.Flags().StringVarP(&extractOutput, "output", "o", "", "Directory to extract files into (required)")
	_ = extractCmd.MarkFlagRequired("output")
	extractCmd.Flags().BoolVar(&extractReadOnly, "read-only", false, "Open package read-only (no write risk)")
}

func runExtract(_ *cobra.Command, args []string) error {
	pathPrefix := ""
	if len(args) > 1 {
		pathPrefix = strings.TrimPrefix(args[1], "/")
	}
	destDir, err := resolveExtractDest(extractOutput)
	if err != nil {
		return err
	}
	ctx := context.Background()
	pkg, err := openPackage(ctx, args[0], extractReadOnly)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()
	files, err := pkg.ListFiles()
	if err != nil {
		return fmt.Errorf("list files: %w", err)
	}
	for _, f := range files {
		if !matchPathPrefix(f.PrimaryPath, pathPrefix) {
			continue
		}
		if err := extractOneFile(ctx, pkg, destDir, f.PrimaryPath); err != nil {
			return err
		}
	}
	return nil
}

func resolveExtractDest(destDir string) (string, error) {
	if destDir == "" {
		return "", fmt.Errorf("output directory is required (use -o or --output)")
	}
	destDir, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("output path: %w", err)
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}
	return destDir, nil
}

func matchPathPrefix(displayPath, pathPrefix string) bool {
	if pathPrefix == "" {
		return true
	}
	return displayPath == pathPrefix || strings.HasPrefix(displayPath+"/", pathPrefix+"/")
}

func extractOneFile(ctx context.Context, pkg novuspack.Package, destDir, displayPath string) error {
	storedPath := "/" + displayPath
	data, err := pkg.ReadFile(ctx, storedPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", storedPath, err)
	}
	destPath := filepath.Join(destDir, filepath.FromSlash(displayPath))
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("create directory for %s: %w", destPath, err)
	}
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", destPath, err)
	}
	return nil
}
