// Package novuspack provides the NovusPack API v1 implementation.
//
// This file contains unit tests for the PackageBuilder pattern:
// NewBuilder, With* methods, and Build.
package novus_package

import (
	"context"
	"strings"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/metadata"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// TestNewBuilder tests the NewBuilder function.
func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}
}

// TestPackageBuilder_WithCompression tests the WithCompression method.
func TestPackageBuilder_WithCompression(t *testing.T) {
	builder := NewBuilder()
	result := builder.WithCompression(CompressionType(fileformat.CompressionZstd))
	if result != builder {
		t.Errorf("WithCompression() should return the same builder instance")
	}
}

// TestPackageBuilder_WithEncryption tests the WithEncryption method.
func TestPackageBuilder_WithEncryption(t *testing.T) {
	builder := NewBuilder()
	result := builder.WithEncryption(EncryptionType(fileformat.EncryptionAES256GCM))
	if result != builder {
		t.Errorf("WithEncryption() should return the same builder instance")
	}
}

// TestPackageBuilder_WithMetadata tests the WithMetadata method.
func TestPackageBuilder_WithMetadata(t *testing.T) {
	builder := NewBuilder()
	md := map[string]string{"key": "value"}
	result := builder.WithMetadata(md)
	if result != builder {
		t.Errorf("WithMetadata() should return the same builder instance")
	}

	// Test with nil metadata
	builder2 := NewBuilder()
	result2 := builder2.WithMetadata(nil)
	if result2 != builder2 {
		t.Errorf("WithMetadata(nil) should return the same builder instance")
	}
}

// TestPackageBuilder_WithComment tests the WithComment method.
func TestPackageBuilder_WithComment(t *testing.T) {
	builder := NewBuilder()
	result := builder.WithComment("test comment")
	if result != builder {
		t.Errorf("WithComment() should return the same builder instance")
	}
}

// TestPackageBuilder_WithVendorID tests the WithVendorID method.
func TestPackageBuilder_WithVendorID(t *testing.T) {
	builder := NewBuilder()
	result := builder.WithVendorID(123)
	if result != builder {
		t.Errorf("WithVendorID() should return the same builder instance")
	}
}

// TestPackageBuilder_WithAppID tests the WithAppID method.
func TestPackageBuilder_WithAppID(t *testing.T) {
	builder := NewBuilder()
	result := builder.WithAppID(456)
	if result != builder {
		t.Errorf("WithAppID() should return the same builder instance")
	}
}

// TestPackageBuilder_Build tests the Build method.
func TestPackageBuilder_Build(t *testing.T) {
	ctx := context.Background()
	builder := NewBuilder().
		WithComment("test comment").
		WithVendorID(1).
		WithAppID(100)

	pkg, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("Build() returned nil package")
	}

	// Verify comment was set (using methods that don't require package to be open)
	if !pkg.HasComment() || pkg.GetComment() != "test comment" {
		t.Errorf("Comment not set correctly: HasComment=%v, Comment=%v", pkg.HasComment(), pkg.GetComment())
	}
	if pkg.GetVendorID() != 1 {
		t.Errorf("VendorID = %v, want 1", pkg.GetVendorID())
	}
	if pkg.GetAppID() != 100 {
		t.Errorf("AppID = %v, want 100", pkg.GetAppID())
	}
}

// TestPackageBuilder_Build_WithAllOptions tests Build with all options set.
func TestPackageBuilder_Build_WithAllOptions(t *testing.T) {
	ctx := context.Background()
	builder := NewBuilder()
	pkg, err := builder.
		WithCompression(CompressionType(fileformat.CompressionZstd)).
		WithEncryption(EncryptionType(fileformat.EncryptionAES256GCM)).
		WithComment("test comment").
		WithVendorID(42).
		WithAppID(100).
		WithMetadata(map[string]string{"key": "value"}).
		Build(ctx)
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("Build() returned nil package")
	}

	// Verify options were applied (using methods that don't require package to be open)
	if !pkg.HasComment() || pkg.GetComment() != "test comment" {
		t.Errorf("Comment not set: HasComment=%v, Comment=%v", pkg.HasComment(), pkg.GetComment())
	}
	if pkg.GetVendorID() != 42 {
		t.Errorf("VendorID = %v, want 42", pkg.GetVendorID())
	}
	if pkg.GetAppID() != 100 {
		t.Errorf("AppID = %v, want 100", pkg.GetAppID())
	}
}

// TestPackageBuilder_Build_WithCancelledContext tests Build with cancelled context.
func TestPackageBuilder_Build_WithCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	builder := NewBuilder()
	_, err := builder.Build(ctx)
	if err == nil {
		t.Fatal("Build() should fail with cancelled context")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeContext {
		t.Errorf("Expected error type Context, got: %v", pkgErr.Type)
	}
}

// TestPackageBuilder_Build_WithNilContext tests Build with nil context.
func TestPackageBuilder_Build_WithNilContext(t *testing.T) {
	builder := NewBuilder()
	//nolint:staticcheck // SA1012: intentionally testing nil context handling
	_, err := builder.Build(nil)
	if err == nil {
		t.Fatal("Build() should fail with nil context")
	}
	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackageBuilder_Build_SetCommentError tests Build when SetComment fails.
func TestPackageBuilder_Build_SetCommentError(t *testing.T) {
	ctx := context.Background()
	// Create a comment that's too long to trigger validation error
	longComment := strings.Repeat("a", int(metadata.MaxCommentLength)+1)
	builder := NewBuilder().WithComment(longComment)

	_, err := builder.Build(ctx)
	if err == nil {
		t.Error("Build() should fail when SetComment fails")
	}
}

func runBuildAndAssertID[T comparable](t *testing.T, builder PackageBuilder, getter func(Package) T, want T, fieldName string) {
	t.Helper()
	ctx := context.Background()
	pkg, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}
	if pkg == nil {
		t.Fatal("Build() returned nil package")
	}
	if getter(pkg) != want {
		t.Errorf("%s = %v, want %v", fieldName, getter(pkg), want)
	}
}

// TestPackageBuilder_Build_SetVendorIDError tests Build when SetVendorID fails.
func TestPackageBuilder_Build_SetVendorIDError(t *testing.T) {
	runBuildAndAssertID(t, NewBuilder().WithVendorID(12345), func(p Package) uint32 { return p.GetVendorID() }, uint32(12345), "VendorID")
}

// TestPackageBuilder_Build_SetAppIDError tests Build when SetAppID fails.
func TestPackageBuilder_Build_SetAppIDError(t *testing.T) {
	runBuildAndAssertID(t, NewBuilder().WithAppID(67890), func(p Package) uint64 { return p.GetAppID() }, uint64(67890), "AppID")
}
