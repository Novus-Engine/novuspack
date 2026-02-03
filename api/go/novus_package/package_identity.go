// This file implements package-level AppID and VendorID management operations.
// It contains methods for setting, getting, clearing, and checking AppID and VendorID
// as specified in api_metadata.md Sections 2-4. This file should contain only
// Package-level identity operations (SetAppID, GetAppID, ClearAppID, HasAppID,
// SetVendorID, GetVendorID, ClearVendorID, HasVendorID, SetPackageIdentity,
// GetPackageIdentity, ClearPackageIdentity).
//
// Specification: api_metadata.md: 1. Comment Management

package novus_package

import (
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// ensureInfoNotNil returns *PackageError if p.Info is nil.
func (p *filePackage) ensureInfoNotNil() error {
	if p.Info == nil {
		return pkgerrors.NewPackageError(pkgerrors.ErrTypeValidation, "package info is nil", nil, struct{}{})
	}
	return nil
}

// SetAppID sets or updates the package AppID.
//
// Sets the AppID in PackageInfo (the single source of truth) and syncs to header.
//
// Parameters:
//   - appID: Application ID to set (uint64)
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 1. Comment Management
func (p *filePackage) SetAppID(appID uint64) error {
	if err := p.ensureInfoNotNil(); err != nil {
		return err
	}
	p.Info.AppID = appID
	p.Info.MetadataVersion++
	return nil
}

// GetAppID retrieves the current package AppID.
//
// Returns the AppID from PackageInfo (the single source of truth). This is a pure data access method
// and does not require context.
//
// Returns:
//   - uint64: Current package AppID
//
// Specification: api_metadata.md: 2. AppID Management
func (p *filePackage) GetAppID() uint64 {
	if p.Info == nil {
		return 0
	}
	return p.Info.AppID
}

// ClearAppID removes the package AppID (sets to 0).
//
// Clears the AppID in both the package header and PackageInfo.
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 2. AppID Management
func (p *filePackage) ClearAppID() error {
	return p.SetAppID(0)
}

// HasAppID checks if the package has an AppID (non-zero).
//
// Returns true if the package has a non-zero AppID, false otherwise.
// This is a pure data access method and does not require context.
//
// Returns:
//   - bool: True if package has an AppID, false otherwise
//
// Specification: api_metadata.md: 2. AppID Management
func (p *filePackage) HasAppID() bool {
	return p.GetAppID() != 0
}

// SetVendorID sets or updates the package VendorID.
//
// Sets the VendorID in PackageInfo (the single source of truth) and syncs to header.
//
// Parameters:
//   - vendorID: Vendor ID to set (uint32)
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 2. AppID Management
func (p *filePackage) SetVendorID(vendorID uint32) error {
	if err := p.ensureInfoNotNil(); err != nil {
		return err
	}
	p.Info.VendorID = vendorID
	p.Info.MetadataVersion++
	return nil
}

// GetVendorID retrieves the current package VendorID.
//
// Returns the VendorID from PackageInfo (the single source of truth). This is a pure data access method
// and does not require context.
//
// Returns:
//   - uint32: Current package VendorID
//
// Specification: api_metadata.md: 3. VendorID Management
func (p *filePackage) GetVendorID() uint32 {
	if p.Info == nil {
		return 0
	}
	return p.Info.VendorID
}

// ClearVendorID removes the package VendorID (sets to 0).
//
// Clears the VendorID in both the package header and PackageInfo.
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 3. VendorID Management
func (p *filePackage) ClearVendorID() error {
	return p.SetVendorID(0)
}

// HasVendorID checks if the package has a VendorID (non-zero).
//
// Returns true if the package has a non-zero VendorID, false otherwise.
// This is a pure data access method and does not require context.
//
// Returns:
//   - bool: True if package has a VendorID, false otherwise
//
// Specification: api_metadata.md: 3. VendorID Management
func (p *filePackage) HasVendorID() bool {
	return p.GetVendorID() != 0
}

// SetPackageIdentity sets both VendorID and AppID.
//
// Sets both the VendorID and AppID in both the package header and PackageInfo
// for consistency.
//
// Parameters:
//   - vendorID: Vendor ID to set (uint32)
//   - appID: Application ID to set (uint64)
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 3. VendorID Management
func (p *filePackage) SetPackageIdentity(vendorID uint32, appID uint64) error {
	// Set VendorID first
	if err := p.SetVendorID(vendorID); err != nil {
		return err
	}

	// Set AppID
	if err := p.SetAppID(appID); err != nil {
		return err
	}

	return nil
}

// GetPackageIdentity gets both VendorID and AppID.
//
// Returns both the VendorID and AppID from the package header. This is a pure
// data access method and does not require context.
//
// Returns:
//   - uint32: Current package VendorID
//   - uint64: Current package AppID
//
// Specification: api_metadata.md: 4. Combined Management
func (p *filePackage) GetPackageIdentity() (vendorID uint32, appID uint64) {
	return p.GetVendorID(), p.GetAppID()
}

// ClearPackageIdentity clears both VendorID and AppID.
//
// Clears both the VendorID and AppID in both the package header and PackageInfo.
//
// Returns:
//   - error: *PackageError on failure
//
// Specification: api_metadata.md: 4. Combined Management
func (p *filePackage) ClearPackageIdentity() error {
	// Clear both
	if err := p.ClearVendorID(); err != nil {
		return err
	}

	if err := p.ClearAppID(); err != nil {
		return err
	}

	return nil
}
