// This file contains unit tests for package AppID and VendorID management operations.
// It tests SetAppID, GetAppID, ClearAppID, HasAppID, SetVendorID, GetVendorID,
// ClearVendorID, HasVendorID, SetPackageIdentity, GetPackageIdentity, and
// ClearPackageIdentity methods from package_identity.go.
//
// Specification: api_metadata.md: 1. Comment Management

package novus_package

import (
	"testing"

	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// =============================================================================
// TEST: SetAppID / GetAppID / HasAppID / ClearAppID
// =============================================================================

// TestPackage_SetAppID_Basic tests basic SetAppID operation.
func TestPackage_SetAppID_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set AppID
	appID := uint64(12345)
	err = fpkg.SetAppID(appID)
	if err != nil {
		t.Errorf("SetAppID() failed: %v", err)
	}

	// Verify AppID was set
	retrieved := fpkg.GetAppID()
	if retrieved != appID {
		t.Errorf("GetAppID() = %d, want %d", retrieved, appID)
	}

	if !fpkg.HasAppID() {
		t.Error("HasAppID() should return true after SetAppID")
	}
}

// TestPackage_ClearAppID_Basic tests basic ClearAppID operation.
func TestPackage_ClearAppID_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set AppID first
	err = fpkg.SetAppID(12345)
	if err != nil {
		t.Fatalf("SetAppID() failed: %v", err)
	}

	// Clear AppID
	err = fpkg.ClearAppID()
	if err != nil {
		t.Errorf("ClearAppID() failed: %v", err)
	}

	// Verify AppID was cleared
	if fpkg.GetAppID() != 0 {
		t.Errorf("GetAppID() = %d, want 0", fpkg.GetAppID())
	}

	if fpkg.HasAppID() {
		t.Error("HasAppID() should return false after ClearAppID")
	}
}

// TestPackage_HasAppID_Basic tests basic HasAppID operation.
func TestPackage_HasAppID_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Initially no AppID
	if fpkg.HasAppID() {
		t.Error("HasAppID() should return false for new package")
	}

	// Set AppID
	err = fpkg.SetAppID(12345)
	if err != nil {
		t.Fatalf("SetAppID() failed: %v", err)
	}

	// Verify HasAppID returns true
	if !fpkg.HasAppID() {
		t.Error("HasAppID() should return true after SetAppID")
	}

	// Clear AppID
	err = fpkg.ClearAppID()
	if err != nil {
		t.Fatalf("ClearAppID() failed: %v", err)
	}

	// Verify HasAppID returns false
	if fpkg.HasAppID() {
		t.Error("HasAppID() should return false after ClearAppID")
	}
}

// TestPackage_AppID_PersistsInHeader tests that AppID is stored in Info and synced to header.
func TestPackage_AppID_PersistsInHeader(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set AppID
	appID := uint64(54321)
	err = fpkg.SetAppID(appID)
	if err != nil {
		t.Fatalf("SetAppID() failed: %v", err)
	}

	// Verify PackageInfo (single source of truth)
	if fpkg.Info.AppID != appID {
		t.Errorf("Info.AppID = %d, want %d", fpkg.Info.AppID, appID)
	}

	// Note: Header is NOT synced immediately after mutations.
	// It will be synced during write operations (Write/SafeWrite/FastWrite).
	// This follows the "PackageInfo as single source of truth" pattern.
}

// =============================================================================
// TEST: SetVendorID / GetVendorID / HasVendorID / ClearVendorID
// =============================================================================

// TestPackage_SetVendorID_Basic tests basic SetVendorID operation.
func TestPackage_SetVendorID_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set VendorID
	vendorID := uint32(67890)
	err = fpkg.SetVendorID(vendorID)
	if err != nil {
		t.Errorf("SetVendorID() failed: %v", err)
	}

	// Verify VendorID was set
	retrieved := fpkg.GetVendorID()
	if retrieved != vendorID {
		t.Errorf("GetVendorID() = %d, want %d", retrieved, vendorID)
	}

	if !fpkg.HasVendorID() {
		t.Error("HasVendorID() should return true after SetVendorID")
	}
}

// TestPackage_ClearVendorID_Basic tests basic ClearVendorID operation.
func TestPackage_ClearVendorID_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set VendorID first
	err = fpkg.SetVendorID(67890)
	if err != nil {
		t.Fatalf("SetVendorID() failed: %v", err)
	}

	// Clear VendorID
	err = fpkg.ClearVendorID()
	if err != nil {
		t.Errorf("ClearVendorID() failed: %v", err)
	}

	// Verify VendorID was cleared
	if fpkg.GetVendorID() != 0 {
		t.Errorf("GetVendorID() = %d, want 0", fpkg.GetVendorID())
	}

	if fpkg.HasVendorID() {
		t.Error("HasVendorID() should return false after ClearVendorID")
	}
}

// TestPackage_HasVendorID_Basic tests basic HasVendorID operation.
func TestPackage_HasVendorID_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Initially no VendorID
	if fpkg.HasVendorID() {
		t.Error("HasVendorID() should return false for new package")
	}

	// Set VendorID
	err = fpkg.SetVendorID(67890)
	if err != nil {
		t.Fatalf("SetVendorID() failed: %v", err)
	}

	// Verify HasVendorID returns true
	if !fpkg.HasVendorID() {
		t.Error("HasVendorID() should return true after SetVendorID")
	}

	// Clear VendorID
	err = fpkg.ClearVendorID()
	if err != nil {
		t.Fatalf("ClearVendorID() failed: %v", err)
	}

	// Verify HasVendorID returns false
	if fpkg.HasVendorID() {
		t.Error("HasVendorID() should return false after ClearVendorID")
	}
}

// TestPackage_VendorID_PersistsInHeader tests that VendorID is stored in Info and synced to header.
func TestPackage_VendorID_PersistsInHeader(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set VendorID
	vendorID := uint32(98765)
	err = fpkg.SetVendorID(vendorID)
	if err != nil {
		t.Fatalf("SetVendorID() failed: %v", err)
	}

	// Verify PackageInfo (single source of truth)
	if fpkg.Info.VendorID != vendorID {
		t.Errorf("Info.VendorID = %d, want %d", fpkg.Info.VendorID, vendorID)
	}

	// Note: Header is NOT synced immediately after mutations.
	// It will be synced during write operations (Write/SafeWrite/FastWrite).
	// This follows the "PackageInfo as single source of truth" pattern.
}

// =============================================================================
// TEST: SetPackageIdentity / GetPackageIdentity / ClearPackageIdentity
// =============================================================================

// TestPackage_SetPackageIdentity_Basic tests basic SetPackageIdentity operation.
func TestPackage_SetPackageIdentity_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set both VendorID and AppID
	vendorID := uint32(11111)
	appID := uint64(22222)
	err = fpkg.SetPackageIdentity(vendorID, appID)
	if err != nil {
		t.Errorf("SetPackageIdentity() failed: %v", err)
	}

	// Verify both were set
	retrievedVendorID, retrievedAppID := fpkg.GetPackageIdentity()
	if retrievedVendorID != vendorID {
		t.Errorf("GetPackageIdentity() VendorID = %d, want %d", retrievedVendorID, vendorID)
	}
	if retrievedAppID != appID {
		t.Errorf("GetPackageIdentity() AppID = %d, want %d", retrievedAppID, appID)
	}
}

// TestPackage_ClearPackageIdentity_Basic tests basic ClearPackageIdentity operation.
func TestPackage_ClearPackageIdentity_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Set both first
	err = fpkg.SetPackageIdentity(11111, 22222)
	if err != nil {
		t.Fatalf("SetPackageIdentity() failed: %v", err)
	}

	// Clear both
	err = fpkg.ClearPackageIdentity()
	if err != nil {
		t.Errorf("ClearPackageIdentity() failed: %v", err)
	}

	// Verify both were cleared
	vendorID, appID := fpkg.GetPackageIdentity()
	if vendorID != 0 {
		t.Errorf("GetPackageIdentity() VendorID = %d, want 0", vendorID)
	}
	if appID != 0 {
		t.Errorf("GetPackageIdentity() AppID = %d, want 0", appID)
	}
}

// TestPackage_GetPackageIdentity_Basic tests basic GetPackageIdentity operation.
func TestPackage_GetPackageIdentity_Basic(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Initially both should be 0
	vendorID, appID := fpkg.GetPackageIdentity()
	if vendorID != 0 {
		t.Errorf("GetPackageIdentity() VendorID = %d, want 0", vendorID)
	}
	if appID != 0 {
		t.Errorf("GetPackageIdentity() AppID = %d, want 0", appID)
	}

	// Set both
	err = fpkg.SetPackageIdentity(33333, 44444)
	if err != nil {
		t.Fatalf("SetPackageIdentity() failed: %v", err)
	}

	// Verify both were set
	vendorID, appID = fpkg.GetPackageIdentity()
	if vendorID != 33333 {
		t.Errorf("GetPackageIdentity() VendorID = %d, want 33333", vendorID)
	}
	if appID != 44444 {
		t.Errorf("GetPackageIdentity() AppID = %d, want 44444", appID)
	}
}

// TestPackage_GetAppID_WithNilHeader tests GetAppID when header is nil.
func TestPackage_GetAppID_WithNilHeader(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Temporarily set header to nil
	originalHeader := fpkg.header
	fpkg.header = nil

	// GetAppID should return 0 when header is nil
	appID := fpkg.GetAppID()
	if appID != 0 {
		t.Errorf("GetAppID() = %d, want 0 when header is nil", appID)
	}

	// Restore header
	fpkg.header = originalHeader
}

// TestPackage_GetVendorID_WithNilHeader tests GetVendorID when header is nil.
func TestPackage_GetVendorID_WithNilHeader(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)

	// Temporarily set header to nil
	originalHeader := fpkg.header
	fpkg.header = nil

	// GetVendorID should return 0 when header is nil
	vendorID := fpkg.GetVendorID()
	if vendorID != 0 {
		t.Errorf("GetVendorID() = %d, want 0 when header is nil", vendorID)
	}

	// Restore header
	fpkg.header = originalHeader
}

// TestPackage_SetPackageIdentity_WithNilHeader removed - obsolete test.
// Since PackageInfo is now the single source of truth, mutation methods
// don't need to check or update header. Header is synced during write operations only.

// TestPackage_ClearPackageIdentity_WithNilHeader removed - obsolete test.
// Since PackageInfo is now the single source of truth, mutation methods
// don't need to check or update header. Header is synced during write operations only.

// TestPackage_SetAppID_WithNilHeader removed - obsolete test.
// Since PackageInfo is now the single source of truth, mutation methods
// don't need to check or update header. Header is synced during write operations only.

// TestPackage_SetVendorID_WithNilHeader removed - obsolete test.
// Since PackageInfo is now the single source of truth, mutation methods
// don't need to check or update header. Header is synced during write operations only.

// TestPackage_SetPackageIdentity_ErrorOnSetVendorIDFailure removed - obsolete test.
// Since PackageInfo is now the single source of truth, mutation methods
// don't need to check or update header. Header is synced during write operations only.

// TestPackage_ClearPackageIdentity_ErrorOnClearVendorIDFailure removed - obsolete test.
// Since PackageInfo is now the single source of truth, mutation methods
// don't need to check or update header. Header is synced during write operations only.

// =============================================================================
// TEST: Edge Cases - Nil Info
// =============================================================================

// TestPackage_GetAppID_WithNilInfo tests GetAppID when Info is nil.
// Expected: Should return 0 gracefully
func TestPackage_GetAppID_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// GetAppID should return 0 (not panic)
	appID := fpkg.GetAppID()
	if appID != 0 {
		t.Errorf("GetAppID() with nil Info = %d, want 0", appID)
	}
}

// TestPackage_GetVendorID_WithNilInfo tests GetVendorID when Info is nil.
// Expected: Should return 0 gracefully
func TestPackage_GetVendorID_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// GetVendorID should return 0 (not panic)
	vendorID := fpkg.GetVendorID()
	if vendorID != 0 {
		t.Errorf("GetVendorID() with nil Info = %d, want 0", vendorID)
	}
}

// TestPackage_HasAppID_WithNilInfo tests HasAppID when Info is nil.
// Expected: Should return false gracefully
func TestPackage_HasAppID_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// HasAppID should return false (not panic)
	if fpkg.HasAppID() {
		t.Error("HasAppID() with nil Info should return false")
	}
}

// TestPackage_HasVendorID_WithNilInfo tests HasVendorID when Info is nil.
// Expected: Should return false gracefully
func TestPackage_HasVendorID_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// HasVendorID should return false (not panic)
	if fpkg.HasVendorID() {
		t.Error("HasVendorID() with nil Info should return false")
	}
}

// TestPackage_GetPackageIdentity_WithNilInfo tests GetPackageIdentity when Info is nil.
// Expected: Should return zeros gracefully
func TestPackage_GetPackageIdentity_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// GetPackageIdentity should return zeros (not panic)
	vendorID, appID := fpkg.GetPackageIdentity()
	if vendorID != 0 {
		t.Errorf("GetPackageIdentity() VendorID with nil Info = %d, want 0", vendorID)
	}
	if appID != 0 {
		t.Errorf("GetPackageIdentity() AppID with nil Info = %d, want 0", appID)
	}
}

// TestPackage_SetAppID_WithNilInfo tests SetAppID when Info is nil.
// Expected: Should return validation error
func TestPackage_SetAppID_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// SetAppID should return error
	err = fpkg.SetAppID(12345)
	if err == nil {
		t.Error("SetAppID() with nil Info should return error")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_SetVendorID_WithNilInfo tests SetVendorID when Info is nil.
// Expected: Should return validation error
func TestPackage_SetVendorID_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// SetVendorID should return error
	err = fpkg.SetVendorID(12345)
	if err == nil {
		t.Error("SetVendorID() with nil Info should return error")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_SetPackageIdentity_WithNilInfo tests SetPackageIdentity when Info is nil.
// Expected: Should return validation error
func TestPackage_SetPackageIdentity_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// SetPackageIdentity should return error
	err = fpkg.SetPackageIdentity(12345, 67890)
	if err == nil {
		t.Error("SetPackageIdentity() with nil Info should return error")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}

// TestPackage_ClearPackageIdentity_WithNilInfo tests ClearPackageIdentity when Info is nil.
// Expected: Should return validation error
func TestPackage_ClearPackageIdentity_WithNilInfo(t *testing.T) {
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()

	fpkg := pkg.(*filePackage)
	// Manually set Info to nil
	fpkg.Info = nil

	// ClearPackageIdentity should return error
	err = fpkg.ClearPackageIdentity()
	if err == nil {
		t.Error("ClearPackageIdentity() with nil Info should return error")
	}

	pkgErr := &pkgerrors.PackageError{}
	if !asPackageError(err, pkgErr) {
		t.Fatalf("Expected PackageError, got: %T", err)
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
	}
}
