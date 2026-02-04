package cmd

import (
	"context"
	"fmt"
	"os"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var headerCmd = &cobra.Command{
	Use:   "header <package path>",
	Short: "Print package header (format version, index offset, etc.)",
	Args:  cobra.ExactArgs(1),
	RunE:  runHeader,
}

func runHeader(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	h, err := novuspack.ReadHeaderFromPath(ctx, path)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Magic:         0x%08X\n", h.Magic)
	_, _ = fmt.Fprintf(os.Stdout, "FormatVersion: %d\n", h.FormatVersion)
	_, _ = fmt.Fprintf(os.Stdout, "IndexStart:    %d\n", h.IndexStart)
	_, _ = fmt.Fprintf(os.Stdout, "Flags:         0x%08X\n", h.Flags)
	return nil
}
