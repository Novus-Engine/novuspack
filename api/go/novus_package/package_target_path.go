// This file implements the SetTargetPath method for the Package interface.
// It provides the ability to change the package's target write path after creation or opening.
//
// Specification: api_basic_operations.md: 1. Context Integration

package novus_package

import (
	"context"
	"os"
	"path/filepath"

	"github.com/novus-engine/novuspack/api/go/internal"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// SetTargetPath changes the package's target write path.
//
// This function changes the target path for an existing package that will be used
// when calling Write, SafeWrite, or FastWrite. This is useful when you want to write
// an existing package (either newly created or opened from disk) to a different location.
//
// Early Validation: This function validates the path and target directory immediately
// (requiring minimal I/O), enabling early error detection before write operations begin.
//
// Path Validation: This function validates that the provided path is valid and the
// target directory is writable, even though it doesn't write to disk. This validation
// requires minimal filesystem I/O to check directory existence and permissions.
//
// Signature Clearing: If the package is signed and the new path differs from the
// current path, this function MUST clear all signature information from the in-memory
// package. This is required because signed packages are immutable and writing to a new
// location creates a new, unsigned package.
//
// Important: Signature clearing only occurs when the new path differs from the current
// path. If SetTargetPath is called with the same path as the current path, signatures
// are NOT cleared.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - path: New file system path where the package will be written
//
// Returns:
//   - error: *PackageError on failure
//
// Error Conditions:
//   - ErrTypeValidation: Invalid or malformed file path, target directory does not exist,
//     target directory is not writable, or insufficient permissions
//   - ErrTypeSecurity: Insufficient permissions to access target directory
//   - ErrTypeContext: Context cancellation or timeout exceeded
//
// Example:
//
//	pkg, _ := OpenPackage(ctx, "source.nvpk")
//	err := pkg.SetTargetPath(ctx, "/new/location/output.nvpk")
//	if err != nil {
//	    return err
//	}
//	err = pkg.Write(ctx) // Writes to /new/location/output.nvpk
//
// Specification: api_basic_operations.md: 1. Context Integration
func (p *filePackage) SetTargetPath(ctx context.Context, path string) error {
	// Check context
	if err := internal.CheckContext(ctx, "SetTargetPath"); err != nil {
		return err
	}

	// Validate path is not empty
	if path == "" {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"target path cannot be empty",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    path,
				Expected: "non-empty file path",
			},
		)
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Get the directory part of the path
	dir := filepath.Dir(cleanPath)

	// Check if directory exists
	dirInfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return pkgerrors.NewPackageError(
				pkgerrors.ErrTypeValidation,
				"target directory does not exist",
				err,
				pkgerrors.ValidationErrorContext{
					Field:    "path",
					Value:    dir,
					Expected: "existing directory",
				},
			)
		}
		return pkgerrors.NewPackageError[struct{}](
			pkgerrors.ErrTypeIO,
			"failed to stat target directory",
			err,
			struct{}{},
		)
	}

	// Verify it's a directory
	if !dirInfo.IsDir() {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"target path parent is not a directory",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "path",
				Value:    dir,
				Expected: "directory",
			},
		)
	}

	// Check if directory is writable by attempting to create a temporary file
	tempFile, err := os.CreateTemp(dir, ".nvpk-write-test-*")
	if err != nil {
		return pkgerrors.NewPackageError[struct{}](
			pkgerrors.ErrTypeSecurity,
			"target directory is not writable",
			err,
			struct{}{},
		)
	}
	// Clean up the temp file
	tempPath := tempFile.Name()
	_ = tempFile.Close()
	_ = os.Remove(tempPath)

	// Note: In v2, if the new path differs from current path and package is signed,
	// signatures would be cleared here. For v1, packages are never signed.

	// Update the target path
	p.FilePath = cleanPath

	return nil
}
