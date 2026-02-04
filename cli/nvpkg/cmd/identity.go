package cmd

import (
	"context"
	"fmt"
	"os"

	novuspack "github.com/novus-engine/novuspack/api/go"
	"github.com/spf13/cobra"
)

var identityCmd = &cobra.Command{
	Use:   "identity <package path> [--vendor-id N] [--app-id N]",
	Short: "Get or set package Vendor ID and App ID",
	Long:  "With no flags, prints Vendor ID and App ID. Use --vendor-id and/or --app-id to set. Changes are written to disk.",
	Args:  cobra.ExactArgs(1),
	RunE:  runIdentity,
}

var (
	identityVendorID uint32
	identityAppID    uint64
)

func init() {
	identityCmd.Flags().Uint32Var(&identityVendorID, "vendor-id", 0, "Set Vendor ID (0 to leave unchanged)")
	identityCmd.Flags().Uint64Var(&identityAppID, "app-id", 0, "Set App ID (0 to leave unchanged)")
}

func runIdentity(_ *cobra.Command, args []string) error {
	path := args[0]
	ctx := context.Background()

	pkg, err := novuspack.OpenPackage(ctx, path)
	if err != nil {
		return fmt.Errorf("open package: %w", err)
	}
	defer func() { _ = pkg.Close() }()

	// Check if we're setting (either flag provided with non-zero or explicit set)
	setVendor := identityVendorID != 0
	setApp := identityAppID != 0
	// Allow setting to 0 via flag: if user passed --vendor-id 0 we still "set" to clear
	// Cobra doesn't distinguish "not set" vs "set to 0" easily; we use non-zero for "set" here.
	// So --vendor-id 0 and --app-id 0 mean "don't change" unless we add a separate --clear-vendor-id.
	if setVendor {
		if err := pkg.SetVendorID(identityVendorID); err != nil {
			return fmt.Errorf("set vendor-id: %w", err)
		}
	}
	if setApp {
		if err := pkg.SetAppID(identityAppID); err != nil {
			return fmt.Errorf("set app-id: %w", err)
		}
	}
	if setVendor || setApp {
		if err := pkg.Write(ctx); err != nil {
			return fmt.Errorf("write: %w", err)
		}
		_, _ = fmt.Fprintf(os.Stdout, "Identity updated\n")
		return nil
	}
	_, _ = fmt.Fprintf(os.Stdout, "Vendor ID: %d\n", pkg.GetVendorID())
	_, _ = fmt.Fprintf(os.Stdout, "App ID:    %d\n", pkg.GetAppID())
	return nil
}
