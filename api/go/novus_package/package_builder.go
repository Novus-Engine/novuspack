// This file implements the PackageBuilder pattern for creating packages with
// complex configurations using a fluent interface. It contains the PackageBuilder
// interface and implementation for building packages with method chaining as
// specified in api_basic_operations.md. This file should contain builder
// methods for configuring compression, encryption, metadata, comments, and
// vendor/app IDs before building the package.
//
// Specification: api_basic_operations.md: 1. Context Integration

// Package novuspack provides the NovusPack API v1 implementation.
//
// This file implements the PackageBuilder pattern for creating packages with
// complex configurations using a fluent interface.
package novus_package

import (
	"context"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// PackageBuilder defines the interface for building packages with a fluent API.
//
// PackageBuilder provides a fluent interface for creating packages with complex
// configurations, improving code readability and reducing parameter complexity.
//
// Specification: api_basic_operations.md: 1. Context Integration
type PackageBuilder interface {
	WithCompression(comp CompressionType) PackageBuilder
	WithEncryption(enc EncryptionType) PackageBuilder
	WithMetadata(metadata map[string]string) PackageBuilder
	WithComment(comment string) PackageBuilder
	WithVendorID(vendorID uint32) PackageBuilder
	WithAppID(appID uint64) PackageBuilder
	Build(ctx context.Context) (Package, error)
}

// packageBuilder is the concrete implementation of PackageBuilder.
type packageBuilder struct {
	compression CompressionType
	encryption  EncryptionType
	metadata    map[string]string
	comment     string
	vendorID    uint32
	appID       uint64
}

// NewBuilder creates a new package builder.
//
// Returns a new PackageBuilder instance with default values that can be
// configured using the fluent interface methods.
//
// Specification: api_basic_operations.md: 1. Context Integration
func NewBuilder() PackageBuilder {
	return &packageBuilder{
		compression: CompressionType(fileformat.CompressionNone),
		encryption:  EncryptionType(fileformat.EncryptionNone),
		metadata:    make(map[string]string),
		comment:     "",
		vendorID:    0,
		appID:       0,
	}
}

// WithCompression sets the compression type for the package.
func (b *packageBuilder) WithCompression(comp CompressionType) PackageBuilder {
	b.compression = comp
	return b
}

// WithEncryption sets the encryption type for the package.
func (b *packageBuilder) WithEncryption(enc EncryptionType) PackageBuilder {
	b.encryption = enc
	return b
}

// WithMetadata sets metadata for the package.
func (b *packageBuilder) WithMetadata(metadata map[string]string) PackageBuilder {
	if metadata != nil {
		b.metadata = metadata
	} else {
		b.metadata = make(map[string]string)
	}
	return b
}

// WithComment sets the package comment.
func (b *packageBuilder) WithComment(comment string) PackageBuilder {
	b.comment = comment
	return b
}

// WithVendorID sets the vendor identifier for the package.
func (b *packageBuilder) WithVendorID(vendorID uint32) PackageBuilder {
	b.vendorID = vendorID
	return b
}

// WithAppID sets the application identifier for the package.
func (b *packageBuilder) WithAppID(appID uint64) PackageBuilder {
	b.appID = appID
	return b
}

// Build creates a new package with the configured options.
//
// This method creates a new package using NewPackage and applies all the
// configured options. The package is created in memory and must be written
// to disk using one of the Write methods.
//
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//
// Returns:
//   - Package: The created package instance
//   - error: Error if package creation fails
func (b *packageBuilder) Build(ctx context.Context) (Package, error) {
	// Validate context
	if ctx == nil {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "context cannot be nil", nil, struct{}{})
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeContext, "context cancelled", ctx.Err(), struct{}{})
	default:
	}

	// Create new package
	pkg, err := NewPackage()
	if err != nil {
		return nil, err
	}

	// Apply builder options
	// Note: Since Create requires a path and we don't have one yet,
	// we'll set the options directly on the package info
	filePkg, ok := pkg.(*filePackage)
	if !ok {
		return nil, pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "unexpected package type from NewPackage", nil, struct{}{})
	}

	// Set comment (Info is the single source of truth)
	if b.comment != "" {
		if err := filePkg.SetComment(b.comment); err != nil {
			return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "failed to set comment")
		}
	}

	// Set vendor ID and app ID using setters (ensures proper sync and validation)
	if b.vendorID != 0 {
		if err := filePkg.SetVendorID(b.vendorID); err != nil {
			return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "failed to set vendor ID")
		}
	}
	if b.appID != 0 {
		if err := filePkg.SetAppID(b.appID); err != nil {
			return nil, pkgerrors.WrapError(err, pkgerrors.ErrTypeValidation, "failed to set app ID")
		}
	}

	// TODO: Apply compression and encryption settings when Write is implemented
	// TODO: Apply metadata when metadata system is implemented

	return pkg, nil
}
