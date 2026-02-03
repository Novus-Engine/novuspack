// This file implements session base management for the Package interface.
// It provides methods for setting, getting, clearing, and checking the session base path,
// which is used for automatic path derivation when adding files from absolute filesystem paths.
//
// Specification: api_basic_operations.md: 19 Package Session Base Management

package novus_package

import (
	"path/filepath"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// SetSessionBase explicitly sets the package-level session base path.
//
// This method allows setting the session base before any file operations.
// The session base controls how absolute filesystem paths are converted to stored package paths.
//
// The session base is automatically established from the first absolute path when
// AddFileOptions.BasePath is not explicitly set. This method allows explicit control
// over the session base.
//
// Parameters:
//   - basePath: The filesystem base directory to use for path derivation (must be an absolute path)
//
// Returns:
//   - error: *PackageError with ErrTypeValidation if the path format is invalid
//
// Example:
//
//	pkg := NewPackage()
//	err := pkg.SetSessionBase("/home/user/project")
//	if err != nil {
//	    return err
//	}
//
// Specification: api_basic_operations.md: 19.4 Package.SetSessionBase Method
func (p *filePackage) SetSessionBase(basePath string) error {
	// Validate that basePath is not empty
	if basePath == "" {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"session base path cannot be empty",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "basePath",
				Value:    basePath,
				Expected: "non-empty absolute path",
			},
		)
	}

	// Validate that basePath is an absolute path
	if !filepath.IsAbs(basePath) {
		return pkgerrors.NewPackageError(
			pkgerrors.ErrTypeValidation,
			"session base path must be an absolute path",
			nil,
			pkgerrors.ValidationErrorContext{
				Field:    "basePath",
				Value:    basePath,
				Expected: "absolute filesystem path",
			},
		)
	}

	// Clean the path to normalize it
	p.sessionBase = filepath.Clean(basePath)
	return nil
}

// GetSessionBase returns the current session base path.
//
// Returns empty string if no session base has been established.
//
// Returns:
//   - string: The current session base path, or empty string if not set
//
// Example:
//
//	basePath := pkg.GetSessionBase()
//	if basePath == "" {
//	    fmt.Println("No session base set")
//	}
//
// Specification: api_basic_operations.md: 19.5 Package.GetSessionBase Method
func (p *filePackage) GetSessionBase() string {
	return p.sessionBase
}

// ClearSessionBase clears the package-level session base.
//
// Subsequent absolute paths will establish a new session base when
// AddFileOptions.BasePath is not explicitly set.
//
// Example:
//
//	pkg.ClearSessionBase()
//	// Next absolute file path will establish new session base
//
// Specification: api_basic_operations.md: 19.6 Package.ClearSessionBase Method
func (p *filePackage) ClearSessionBase() {
	p.sessionBase = ""
}

// HasSessionBase returns true if a session base is currently set.
//
// Returns:
//   - bool: true if session base is set, false otherwise
//
// Example:
//
//	if pkg.HasSessionBase() {
//	    fmt.Printf("Session base: %s\n", pkg.GetSessionBase())
//	}
//
// Specification: api_basic_operations.md: 19.7 Package.HasSessionBase Method
func (p *filePackage) HasSessionBase() bool {
	return p.sessionBase != ""
}
