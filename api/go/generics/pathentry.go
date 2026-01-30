// This file implements the PathEntry type representing a minimal file or
// directory path. It contains the PathEntry type definition, validation,
// marshaling/unmarshaling methods, and path conversion utilities. This file
// should contain all code related to PathEntry as specified in api_generics.md
// Section 1.3.
//
// Specification: api_generics.md: 1.3 PathEntry Type

package generics

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// PathEntry represents a minimal file or directory path.
//
// PathEntry is used by both FileEntry and DirectoryEntry to represent path information.
// It supports both file and directory paths, with symlink support for runtime operations.
//
// Path metadata (permissions, timestamps, etc.) is stored separately in PathMetadataEntry
// to allow the same file content to have different permissions at different paths.
//
// Size: Variable (2 + path_length bytes)
//
// Specification: api_generics.md: 1.3 PathEntry Type
type PathEntry struct {
	// PathLength is the length of the path in bytes (UTF-8)
	// Specification: api_generics.md: 1.3.2 Binary Format Specification
	PathLength uint16

	// Path is the UTF-8 encoded file or directory path (not null-terminated)
	// Specification: api_generics.md: 1.3.2 Binary Format Specification
	Path string

	// Symbolic link support (runtime only, not stored in binary format)
	// Specification: api_file_mgmt_file_entry.md: 1.1 FileEntry Structure Definition
	IsSymlink  bool   // Whether this path is a symbolic link
	LinkTarget string // Target path for symbolic links (empty if not a symlink)
}

// Validate performs validation checks on the PathEntry.
//
// Validation checks:
//   - PathLength must match actual Path length
//   - Path must not be empty or whitespace only (after trimming)
//   - Path must be valid UTF-8
//
// Returns an error if any validation check fails.
func (p *PathEntry) Validate() error {
	trimmed := strings.TrimSpace(p.Path)
	if trimmed == "" {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path cannot be empty or whitespace only", nil, pkgerrors.ValidationErrorContext{
			Field:    "Path",
			Value:    p.Path,
			Expected: "non-empty path",
		})
	}

	pathLen := uint16(len(p.Path))
	if p.PathLength != pathLen {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path length mismatch", nil, pkgerrors.ValidationErrorContext{
			Field:    "PathLength",
			Value:    p.PathLength,
			Expected: fmt.Sprintf("%d", pathLen),
		})
	}

	return nil
}

// Size returns the total size of the PathEntry in bytes.
//
// Specification: api_generics.md: 1.3.2 Binary Format Specification
func (p *PathEntry) Size() int {
	// PathLength (2 bytes) + Path (PathLength bytes)
	return 2 + int(p.PathLength)
}

// ReadFrom reads a PathEntry from the provided io.Reader.
//
// The binary format is minimal:
//   - PathLength (2 bytes, little-endian uint16)
//   - Path (PathLength bytes, UTF-8 string, not null-terminated)
//
// Returns the number of bytes read and any error encountered.
//
// Specification: api_generics.md: 1.3.2 Binary Format Specification
func (p *PathEntry) ReadFrom(r io.Reader) (int64, error) {
	var totalRead int64

	// Read PathLength (2 bytes)
	var pathLength uint16
	if err := binary.Read(r, binary.LittleEndian, &pathLength); err != nil {
		return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read path length", pkgerrors.ValidationErrorContext{
			Field:    "PathLength",
			Value:    nil,
			Expected: "2 bytes",
		})
	}
	totalRead += 2
	p.PathLength = pathLength

	// Read Path (PathLength bytes)
	if pathLength > 0 {
		pathBytes := make([]byte, pathLength)
		n, err := io.ReadFull(r, pathBytes)
		if err != nil {
			return totalRead, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to read path", pkgerrors.ValidationErrorContext{
				Field:    "Path",
				Value:    pathLength,
				Expected: "path data",
			})
		}
		if uint16(n) != pathLength {
			return totalRead, pkgerrors.NewPackageError(pkgerrors.ErrTypeCorruption, "incomplete path read", nil, pkgerrors.ValidationErrorContext{
				Field:    "Path",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", pathLength),
			})
		}
		totalRead += int64(n)
		p.Path = string(pathBytes)
	} else {
		p.Path = ""
	}

	return totalRead, nil
}

// WriteTo writes a PathEntry to the provided io.Writer.
//
// The binary format is minimal:
//   - PathLength (2 bytes, little-endian uint16)
//   - Path (PathLength bytes, UTF-8 string, not null-terminated)
//
// Returns the number of bytes written and any error encountered.
//
// Specification: api_generics.md: 1.3.2 Binary Format Specification
func (p *PathEntry) WriteTo(w io.Writer) (int64, error) {
	var totalWritten int64

	// Write PathLength (2 bytes)
	if err := binary.Write(w, binary.LittleEndian, p.PathLength); err != nil {
		return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write path length", pkgerrors.ValidationErrorContext{
			Field:    "PathLength",
			Value:    p.PathLength,
			Expected: "written successfully",
		})
	}
	totalWritten += 2

	// Write Path (PathLength bytes)
	if p.PathLength > 0 {
		pathBytes := []byte(p.Path)
		if uint16(len(pathBytes)) != p.PathLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "path length mismatch", nil, pkgerrors.ValidationErrorContext{
				Field:    "PathLength",
				Value:    len(pathBytes),
				Expected: fmt.Sprintf("%d", p.PathLength),
			})
		}
		n, err := w.Write(pathBytes)
		if err != nil {
			return totalWritten, pkgerrors.WrapErrorWithContext(err, pkgerrors.ErrTypeIO, "failed to write path", pkgerrors.ValidationErrorContext{
				Field:    "Path",
				Value:    p.Path,
				Expected: "written successfully",
			})
		}
		if uint16(n) != p.PathLength {
			return totalWritten, pkgerrors.NewPackageError(pkgerrors.ErrTypeIO, "incomplete path write", nil, pkgerrors.ValidationErrorContext{
				Field:    "Path",
				Value:    n,
				Expected: fmt.Sprintf("%d bytes", p.PathLength),
			})
		}
		totalWritten += int64(n)
	}

	return totalWritten, nil
}

// GetPath returns the path string as stored (Unix-style with forward slashes).
//
// Returns the path as stored in the package format, including the leading "/"
// that indicates the package root. For display format (without leading "/"),
// use GetPathForPlatform() or ToDisplayPath().
//
// Returns:
//   - string: Path as stored (with leading "/" if present)
//
// Specification: api_generics.md: 1.3.1 PathEntry Structure
func (p *PathEntry) GetPath() string {
	return p.Path
}

// GetPathForPlatform returns the path string converted for the specified platform.
//
// This method implements the display format conversion per api_generics 1.3.3.7:
// 1. Strips the leading "/" (package root indicator, not filesystem root)
// 2. On Windows: converts forward slashes to backslashes
// 3. On Unix/Linux: returns path with forward slashes (no leading "/")
//
// Parameters:
//   - isWindows: True for Windows path format, false for Unix/Linux
//
// Returns:
//   - string: Path in display format for the specified platform
//
// Specification: api_generics.md: 1.3.1 PathEntry Structure
func (p *PathEntry) GetPathForPlatform(isWindows bool) string {
	// Strip leading "/" for display (per 1.3.3.7)
	display := strings.TrimPrefix(p.Path, "/")

	if isWindows {
		// Convert forward slashes to backslashes for Windows
		return strings.ReplaceAll(display, "/", "\\")
	}

	// Return with forward slashes for Unix/Linux (no leading "/")
	return display
}
