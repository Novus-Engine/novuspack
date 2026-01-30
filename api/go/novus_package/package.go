// This file defines the core Package interface and filePackage implementation.
// It contains Package, PackageReader, and PackageWriter interfaces as specified
// in api_core.md, along with the filePackage struct that implements these interfaces.
// This file should contain the main Package type definition, interface declarations,
// and the NewPackage constructor function.
//
// Specification: api_core.md: 1. Core Interfaces

// Package novuspack provides the NovusPack API v1 implementation.
//
// This package implements the core Package type and lifecycle operations
// for NovusPack file format management. It provides high-level operations
// for creating, opening, closing, validating, and retrieving information
// from NovusPack package files.
//
// Thread Safety:
//   - Package instances are not thread-safe and should not be shared
//     across goroutines without external synchronization.
//   - Use separate Package instances for concurrent operations.
//
// Context Integration:
//   - All I/O methods accept context.Context for cancellation and timeout.
//   - Methods check context cancellation before performing operations.
//   - Context errors are wrapped in PackageError with ErrTypeContext.
//
// Error Handling:
//   - All errors are wrapped in PackageError for structured error handling.
//   - Use type assertions or pkgerrors.As() to access the PackageError type.
//   - Error types include: Validation, IO, Security, Unsupported, Context, Corruption.
//
// Specification: api_basic_operations.md: 1. Context Integration
package novus_package

import (
	"context"
	"os"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/metadata"
)

// =============================================================================
// INTERFACES
// =============================================================================

// PackageReader defines the interface for reading operations on a package.
//
// PackageReader provides methods for reading files, listing files, retrieving
// metadata, and validating package contents.
//
// Specification: api_core.md: 1.1 PackageReader Interface
type PackageReader interface {
	ReadFile(ctx context.Context, path string) ([]byte, error)
	ListFiles() ([]FileInfo, error)
	GetMetadata() (*metadata.PackageMetadata, error)
	Validate(ctx context.Context) error
	GetInfo() (*metadata.PackageInfo, error)
}

// PackageWriter defines the interface for writing operations on a package.
//
// PackageWriter provides methods for adding files, removing files, and writing
// the package to disk. Compression and signing options are configured via package
// state rather than method parameters.
//
// Specification: api_core.md: 1.2 PackageWriter Interface
type PackageWriter interface {
	Write(ctx context.Context) error
	SafeWrite(ctx context.Context, overwrite bool) error
	FastWrite(ctx context.Context) error
}

// Package defines the main interface for NovusPack package operations.
//
// Package combines PackageReader and PackageWriter interfaces, providing
// complete package lifecycle management including opening, closing, and
// defragmentation operations.
//
// Specification: api_core.md: 1.3 Package Interface
type Package interface {
	PackageReader
	PackageWriter
	Create(ctx context.Context, path string) error
	CreateWithOptions(ctx context.Context, path string, options *CreateOptions) error
	Close() error
	CloseWithCleanup(ctx context.Context) error
	IsOpen() bool
	IsReadOnly() bool
	GetPath() string
	Defragment(ctx context.Context) error

	// Target path management
	// Specification: api_core.md: 1.2 PackageWriter Interface
	SetTargetPath(ctx context.Context, path string) error

	// Session base management
	// Specification: api_basic_operations.md: 3.1 Package Implementation Structure
	SetSessionBase(basePath string) error
	GetSessionBase() string
	ClearSessionBase()
	HasSessionBase() bool

	// File management operations
	// Specification: api_basic_operations.md: 3.1 Package Implementation Structure
	AddFile(ctx context.Context, filesystemPath string, options *AddFileOptions) (*metadata.FileEntry, error)
	AddFileFromMemory(ctx context.Context, path string, data []byte, options *AddFileOptions) (*metadata.FileEntry, error)
	AddFilePattern(ctx context.Context, pattern string, options *AddFileOptions) ([]*metadata.FileEntry, error)
	AddDirectory(ctx context.Context, dirPath string, options *AddFileOptions) ([]*metadata.FileEntry, error)

	// File removal operations
	// Specification: api_file_mgmt_removal.md: 2. RemoveFile Package Method
	RemoveFile(ctx context.Context, path string) error
	RemoveFilePattern(ctx context.Context, pattern string) error
	RemoveDirectory(ctx context.Context, dirPath string) error

	// Comment management operations
	// Specification: api_metadata.md: 1. Comment Management
	SetComment(comment string) error
	GetComment() string
	ClearComment() error
	HasComment() bool

	// AppID/VendorID management operations
	// Specification: api_metadata.md: 1. Comment Management
	SetAppID(appID uint64) error
	GetAppID() uint64
	ClearAppID() error
	HasAppID() bool
	SetVendorID(vendorID uint32) error
	GetVendorID() uint32
	ClearVendorID() error
	HasVendorID() bool
	SetPackageIdentity(vendorID uint32, appID uint64) error
	GetPackageIdentity() (uint32, uint64)
	ClearPackageIdentity() error
}

// =============================================================================
// PACKAGE TYPE
// =============================================================================

// filePackage is the concrete implementation of the Package interface.
//
// filePackage provides the main implementation for interacting with NovusPack files.
// It implements Package, PackageReader, and PackageWriter interfaces.
//
// Lifecycle States:
//   - New: Created via NewPackage(), not yet associated with a file
//   - Created: After Create() is called, ready for file operations
//   - Open: After OpenPackage() is called, file is open for reading
//   - Closed: After Close() is called, resources released, operations not allowed
//
// State Transitions:
//   - NewPackage() → Create() → Close()
//   - OpenPackage() → Close()
//   - GetInfo() can be called on New, Created, or Open states
//   - Validate() can only be called on Open state
//
// Memory Management:
//   - FileEntries are loaded on demand to minimize memory usage
//   - Call Close() to release file handles and resources
//   - Package instances should not be shared across goroutines
//
// Specification: api_basic_operations.md: 1. Context Integration
type filePackage struct {
	// Info contains package metadata (version, counts, timestamps)
	Info *metadata.PackageInfo

	// FileEntries contains all file entries (loaded on demand from index)
	FileEntries []*metadata.FileEntry

	// PathMetadataEntries contains path metadata (loaded from special files)
	PathMetadataEntries []*metadata.PathMetadataEntry

	// SpecialFiles maps special file IDs to their entries (metadata, signatures, etc.)
	SpecialFiles map[uint16]*metadata.FileEntry

	// FilePath is the path to the package file on disk
	FilePath string

	// Internal state (unexported fields)
	header      *fileformat.PackageHeader // Binary header structure
	index       *fileformat.FileIndex     // File index structure
	fileHandle  *os.File                  // Open file handle (nil when closed)
	isOpen      bool                      // True when file is open for reading, false when closed
	sessionBase string                    // Package-level session base path for automatic path derivation (runtime only)
}

// =============================================================================
// CONSTRUCTOR
// =============================================================================

// NewPackage creates a new in-memory Package instance.
//
// This constructor creates a package in memory only and does not write to disk.
// The returned Package is in the "New" state and must be configured using Create()
// before it can be used for file operations.
//
// Use Cases:
//   - Create a new package: NewPackage() → Create()
//   - Alternative: Use OpenPackage() to open an existing package
//
// Returns:
//   - Package: A new Package instance in the "New" state
//   - error: Always returns nil for the error (reserved for future use)
//
// Example:
//
//	pkg, err := novuspack.NewPackage()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	err = pkg.Create(ctx, "mypackage.nvpk")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer pkg.Close()
//
// Specification: api_basic_operations.md: 6.1 Package Constructor
func NewPackage() (Package, error) {
	// Initialize package with default values
	pkg := &filePackage{
		FileEntries:  make([]*metadata.FileEntry, 0),
		SpecialFiles: make(map[uint16]*metadata.FileEntry),
		header:       fileformat.NewPackageHeader(),
		isOpen:       false,
	}

	// Initialize PackageInfo with default values
	pkg.Info = metadata.NewPackageInfo()

	return pkg, nil
}
