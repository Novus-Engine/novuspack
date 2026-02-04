package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <package path>",
	Short: "Create a new empty NovusPack package",
	Args:  cobra.ExactArgs(1),
	RunE:  runCreate,
}

var (
	createComment  string
	createVendorID uint32
	createAppID    uint64
	createModeStr  string
)

func init() {
	createCmd.Flags().StringVar(&createComment, "comment", "", "Package comment")
	createCmd.Flags().Uint32Var(&createVendorID, "vendor-id", 0, "Vendor ID")
	createCmd.Flags().Uint64Var(&createAppID, "app-id", 0, "Application ID")
	createCmd.Flags().StringVar(&createModeStr, "mode", "", "File mode for created package file (e.g. 0644)")
}

func runCreate(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	pkg, err := novuspack.NewPackage()
	if err != nil {
		return fmt.Errorf("new package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	var opts *novuspack.CreateOptions
	if createComment != "" || createVendorID != 0 || createAppID != 0 || createModeStr != "" {
		opts = &novuspack.CreateOptions{
			Comment:  createComment,
			VendorID: createVendorID,
			AppID:    createAppID,
		}
		if createModeStr != "" {
			mode, err := parseFileMode(createModeStr)
			if err != nil {
				return fmt.Errorf("mode: %w", err)
			}
			opts.Permissions = mode
		}
	}
	if opts != nil {
		if err := pkg.CreateWithOptions(ctx, path, opts); err != nil {
			return fmt.Errorf("create with options: %w", err)
		}
	} else {
		if err := pkg.Create(ctx, path); err != nil {
			return fmt.Errorf("create: %w", err)
		}
	}
	if err := pkg.Write(ctx); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Created %s\n", path)
	return nil
}

// parseFileMode parses an octal or decimal file mode string (e.g. 0644, 420).
func parseFileMode(s string) (os.FileMode, error) {
	if s == "" {
		return 0, fmt.Errorf("empty mode")
	}
	// Try octal first (0644)
	n, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		// Try decimal
		n, err = strconv.ParseUint(s, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid mode %q", s)
		}
	}
	return os.FileMode(n), nil
}
