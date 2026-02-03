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

func mustNewFilePackage(t *testing.T) (Package, *filePackage) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	return pkg, pkg.(*filePackage)
}

// identityOps holds set/get/has/clear/info callbacks for AppID or VendorID.
// Used by table-driven identity tests to avoid duplicating test logic.
type (
	setFn   func(*filePackage, interface{}) error
	getFn   func(*filePackage) interface{}
	hasFn   func(*filePackage) bool
	clearFn func(*filePackage) error
	infoFn  func(*filePackage) interface{}
)
type identityOps struct {
	set   setFn
	get   getFn
	has   hasFn
	clear clearFn
	info  infoFn
}

func runIdentitySetGetHas(t *testing.T, fpkg *filePackage, val interface{}, ops identityOps, fieldName string) {
	t.Helper()
	if err := ops.set(fpkg, val); err != nil {
		t.Errorf("Set%s() failed: %v", fieldName, err)
	}
	if got := ops.get(fpkg); got != val {
		t.Errorf("Get%s() = %v, want %v", fieldName, got, val)
	}
	if !ops.has(fpkg) {
		t.Errorf("Has%s() should return true after Set", fieldName)
	}
}

func runIdentityClear(t *testing.T, fpkg *filePackage, setVal interface{}, ops identityOps, fieldName string) {
	t.Helper()
	if err := ops.set(fpkg, setVal); err != nil {
		t.Fatalf("Set%s() failed: %v", fieldName, err)
	}
	if err := ops.clear(fpkg); err != nil {
		t.Errorf("Clear%s() failed: %v", fieldName, err)
	}
	if ops.get(fpkg) != uint64(0) && ops.get(fpkg) != uint32(0) {
		t.Errorf("Get%s() after Clear = %v, want 0", fieldName, ops.get(fpkg))
	}
	if ops.has(fpkg) {
		t.Errorf("Has%s() should return false after Clear", fieldName)
	}
}

func runIdentityHasLifecycle(t *testing.T, fpkg *filePackage, setVal interface{}, ops identityOps, fieldName string) {
	t.Helper()
	if ops.has(fpkg) {
		t.Errorf("Has%s() should return false for new package", fieldName)
	}
	if err := ops.set(fpkg, setVal); err != nil {
		t.Fatalf("Set%s() failed: %v", fieldName, err)
	}
	if !ops.has(fpkg) {
		t.Errorf("Has%s() should return true after Set", fieldName)
	}
	if err := ops.clear(fpkg); err != nil {
		t.Fatalf("Clear%s() failed: %v", fieldName, err)
	}
	if ops.has(fpkg) {
		t.Errorf("Has%s() should return false after Clear", fieldName)
	}
}

func runIdentityPersistsInInfo(t *testing.T, fpkg *filePackage, val interface{}, ops identityOps, fieldName string) {
	t.Helper()
	if err := ops.set(fpkg, val); err != nil {
		t.Fatalf("Set%s() failed: %v", fieldName, err)
	}
	if ops.info(fpkg) != val {
		t.Errorf("Info.%s = %v, want %v", fieldName, ops.info(fpkg), val)
	}
}

func makeIdentityOps(set setFn, get getFn, has hasFn, clearOp clearFn, info infoFn) identityOps {
	return identityOps{set: set, get: get, has: has, clear: clearOp, info: info}
}

func appIDIdentityOps() identityOps {
	return makeIdentityOps(
		func(fpkg *filePackage, v interface{}) error { return fpkg.SetAppID(v.(uint64)) },
		func(fpkg *filePackage) interface{} { return fpkg.GetAppID() },
		func(fpkg *filePackage) bool { return fpkg.HasAppID() },
		func(fpkg *filePackage) error { return fpkg.ClearAppID() },
		func(fpkg *filePackage) interface{} { return fpkg.Info.AppID },
	)
}

func vendorIDIdentityOps() identityOps {
	set := func(fpkg *filePackage, v interface{}) error { return fpkg.SetVendorID(v.(uint32)) }
	get := func(fpkg *filePackage) interface{} { return fpkg.GetVendorID() }
	has := func(fpkg *filePackage) bool { return fpkg.HasVendorID() }
	clearFn := func(fpkg *filePackage) error { return fpkg.ClearVendorID() }
	info := func(fpkg *filePackage) interface{} { return fpkg.Info.VendorID }
	return identityOps{set: set, get: get, has: has, clear: clearFn, info: info}
}

var appIDOps = appIDIdentityOps()
var vendorIDOps = vendorIDIdentityOps()

// =============================================================================
// TEST: SetAppID / GetAppID / HasAppID / ClearAppID
// =============================================================================

// TestPackage_SetAppID_Basic tests basic SetAppID operation.
func TestPackage_SetAppID_Basic(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentitySetGetHas(t, fpkg, uint64(12345), appIDOps, "AppID")
}

// TestPackage_ClearAppID_Basic tests basic ClearAppID operation.
func TestPackage_ClearAppID_Basic(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentityClear(t, fpkg, uint64(12345), appIDOps, "AppID")
}

// TestPackage_HasAppID_Basic tests basic HasAppID operation.
func TestPackage_HasAppID_Basic(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentityHasLifecycle(t, fpkg, uint64(12345), appIDOps, "AppID")
}

// TestPackage_AppID_PersistsInHeader tests that AppID is stored in Info and synced to header.
func TestPackage_AppID_PersistsInHeader(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentityPersistsInInfo(t, fpkg, uint64(54321), appIDOps, "AppID")
}

// =============================================================================
// TEST: SetVendorID / GetVendorID / HasVendorID / ClearVendorID
// =============================================================================

// TestPackage_SetVendorID_Basic tests basic SetVendorID operation.
func TestPackage_SetVendorID_Basic(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentitySetGetHas(t, fpkg, uint32(67890), vendorIDOps, "VendorID")
}

// TestPackage_ClearVendorID_Basic tests basic ClearVendorID operation.
func TestPackage_ClearVendorID_Basic(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentityClear(t, fpkg, uint32(67890), vendorIDOps, "VendorID")
}

// TestPackage_HasVendorID_Basic tests basic HasVendorID operation.
func TestPackage_HasVendorID_Basic(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentityHasLifecycle(t, fpkg, uint32(67890), vendorIDOps, "VendorID")
}

// TestPackage_VendorID_PersistsInHeader tests that VendorID is stored in Info and synced to header.
func TestPackage_VendorID_PersistsInHeader(t *testing.T) {
	pkg, fpkg := mustNewFilePackage(t)
	defer func() { _ = pkg.Close() }()
	runIdentityPersistsInInfo(t, fpkg, uint32(98765), vendorIDOps, "VendorID")
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

func runGetterWithNilHeader(t *testing.T, getter func(*filePackage) interface{}, want interface{}, name string) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	fpkg := pkg.(*filePackage)
	originalHeader := fpkg.header
	fpkg.header = nil
	defer func() { fpkg.header = originalHeader }()
	got := getter(fpkg)
	if got != want {
		t.Errorf("%s() = %v, want %v when header is nil", name, got, want)
	}
}

// TestPackage_GetAppID_WithNilHeader tests GetAppID when header is nil.
func TestPackage_GetAppID_WithNilHeader(t *testing.T) {
	runGetterWithNilHeader(t, func(fpkg *filePackage) interface{} { return fpkg.GetAppID() }, uint64(0), "GetAppID")
}

// TestPackage_GetVendorID_WithNilHeader tests GetVendorID when header is nil.
func TestPackage_GetVendorID_WithNilHeader(t *testing.T) {
	runGetterWithNilHeader(t, func(fpkg *filePackage) interface{} { return fpkg.GetVendorID() }, uint32(0), "GetVendorID")
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
	runWithNilInfo(t, func(t *testing.T, fpkg *filePackage) {
		appID := fpkg.GetAppID()
		if appID != 0 {
			t.Errorf("GetAppID() with nil Info = %d, want 0", appID)
		}
	})
}

// TestPackage_GetVendorID_WithNilInfo tests GetVendorID when Info is nil.
// Expected: Should return 0 gracefully
func TestPackage_GetVendorID_WithNilInfo(t *testing.T) {
	runWithNilInfo(t, func(t *testing.T, fpkg *filePackage) {
		vendorID := fpkg.GetVendorID()
		if vendorID != 0 {
			t.Errorf("GetVendorID() with nil Info = %d, want 0", vendorID)
		}
	})
}

func runWithNilInfo(t *testing.T, check func(t *testing.T, fpkg *filePackage)) {
	t.Helper()
	pkg, err := NewPackage()
	if err != nil {
		t.Fatalf("NewPackage() failed: %v", err)
	}
	defer func() { _ = pkg.Close() }()
	fpkg := pkg.(*filePackage)
	fpkg.Info = nil
	check(t, fpkg)
}

// TestPackage_HasAppID_WithNilInfo tests HasAppID when Info is nil.
// Expected: Should return false gracefully
func TestPackage_HasAppID_WithNilInfo(t *testing.T) {
	runWithNilInfo(t, func(t *testing.T, fpkg *filePackage) {
		if fpkg.HasAppID() {
			t.Error("HasAppID() with nil Info should return false")
		}
	})
}

// TestPackage_HasVendorID_WithNilInfo tests HasVendorID when Info is nil.
// Expected: Should return false gracefully
func TestPackage_HasVendorID_WithNilInfo(t *testing.T) {
	runWithNilInfo(t, func(t *testing.T, fpkg *filePackage) {
		if fpkg.HasVendorID() {
			t.Error("HasVendorID() with nil Info should return false")
		}
	})
}

// TestPackage_GetPackageIdentity_WithNilInfo tests GetPackageIdentity when Info is nil.
// Expected: Should return zeros gracefully
func TestPackage_GetPackageIdentity_WithNilInfo(t *testing.T) {
	runWithNilInfo(t, func(t *testing.T, fpkg *filePackage) {
		vendorID, appID := fpkg.GetPackageIdentity()
		if vendorID != 0 {
			t.Errorf("GetPackageIdentity() VendorID with nil Info = %d, want 0", vendorID)
		}
		if appID != 0 {
			t.Errorf("GetPackageIdentity() AppID with nil Info = %d, want 0", appID)
		}
	})
}

func assertSetWithNilInfoReturnsValidationError(t *testing.T, setFn func(*filePackage) error) {
	t.Helper()
	runWithNilInfo(t, func(t *testing.T, fpkg *filePackage) {
		err := setFn(fpkg)
		if err == nil {
			t.Error("Set with nil Info should return error")
		}
		pkgErr := &pkgerrors.PackageError{}
		if !asPackageError(err, pkgErr) {
			t.Fatalf("Expected PackageError, got: %T", err)
		}
		if pkgErr.Type != pkgerrors.ErrTypeValidation {
			t.Errorf("Expected error type Validation, got: %v", pkgErr.Type)
		}
	})
}

// TestPackage_SetAppID_WithNilInfo tests SetAppID when Info is nil.
// Expected: Should return validation error
func TestPackage_SetAppID_WithNilInfo(t *testing.T) {
	assertSetWithNilInfoReturnsValidationError(t, func(fpkg *filePackage) error {
		return fpkg.SetAppID(12345)
	})
}

// TestPackage_SetVendorID_WithNilInfo tests SetVendorID when Info is nil.
// Expected: Should return validation error
func TestPackage_SetVendorID_WithNilInfo(t *testing.T) {
	assertSetWithNilInfoReturnsValidationError(t, func(fpkg *filePackage) error {
		return fpkg.SetVendorID(12345)
	})
}

// TestPackage_SetPackageIdentity_WithNilInfo tests SetPackageIdentity when Info is nil.
// Expected: Should return validation error
func TestPackage_SetPackageIdentity_WithNilInfo(t *testing.T) {
	assertSetWithNilInfoReturnsValidationError(t, func(fpkg *filePackage) error {
		return fpkg.SetPackageIdentity(12345, 67890)
	})
}

// TestPackage_ClearPackageIdentity_WithNilInfo tests ClearPackageIdentity when Info is nil.
// Expected: Should return validation error
func TestPackage_ClearPackageIdentity_WithNilInfo(t *testing.T) {
	assertSetWithNilInfoReturnsValidationError(t, func(fpkg *filePackage) error {
		return fpkg.ClearPackageIdentity()
	})
}
