// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This file contains unit tests for PackageMetadata structure.
package metadata

import (
	"testing"
)

// TestNewPackageMetadata tests PackageMetadata structure creation.
func TestNewPackageMetadata(t *testing.T) {
	pm := NewPackageMetadata()
	if pm == nil {
		t.Fatal("NewPackageMetadata() returned nil")
	}
	if pm.PackageInfo == nil {
		t.Error("PackageMetadata.PackageInfo should not be nil")
	}
	if pm.FileEntries == nil {
		t.Error("PackageMetadata.FileEntries should not be nil")
	}
	if pm.PathMetadataEntries == nil {
		t.Error("PackageMetadata.PathMetadataEntries should not be nil")
	}
	if pm.SpecialFiles == nil {
		t.Error("PackageMetadata.SpecialFiles should not be nil")
	}
}

// TestPackageMetadata_ContainsPackageInfo tests that PackageMetadata contains PackageInfo fields.
func TestPackageMetadata_ContainsPackageInfo(t *testing.T) {
	pm := NewPackageMetadata()
	if pm.PackageInfo == nil {
		t.Fatal("PackageInfo is nil")
	}
	// Verify PackageInfo fields are accessible
	if pm.FileCount != 0 {
		t.Errorf("FileCount = %v, want 0", pm.FileCount)
	}
	if pm.VendorID != 0 {
		t.Errorf("VendorID = %v, want 0", pm.VendorID)
	}
	if pm.AppID != 0 {
		t.Errorf("AppID = %v, want 0", pm.AppID)
	}
}

// TestPackageMetadata_ContainsFileEntries tests that PackageMetadata contains FileEntries slice.
func TestPackageMetadata_ContainsFileEntries(t *testing.T) {
	pm := NewPackageMetadata()
	if pm.FileEntries == nil {
		t.Fatal("FileEntries is nil")
	}
	if len(pm.FileEntries) != 0 {
		t.Errorf("FileEntries length = %v, want 0", len(pm.FileEntries))
	}
	// Test that we can append to FileEntries
	pm.FileEntries = append(pm.FileEntries, nil)
	if len(pm.FileEntries) != 1 {
		t.Errorf("FileEntries length after append = %v, want 1", len(pm.FileEntries))
	}
}

// TestPackageMetadata_ContainsPathMetadataEntries tests that PackageMetadata contains PathMetadataEntries slice.
func TestPackageMetadata_ContainsPathMetadataEntries(t *testing.T) {
	pm := NewPackageMetadata()
	if pm.PathMetadataEntries == nil {
		t.Fatal("PathMetadataEntries is nil")
	}
	if len(pm.PathMetadataEntries) != 0 {
		t.Errorf("PathMetadataEntries length = %v, want 0", len(pm.PathMetadataEntries))
	}
	// Test that we can append to PathMetadataEntries
	pm.PathMetadataEntries = append(pm.PathMetadataEntries, nil)
	if len(pm.PathMetadataEntries) != 1 {
		t.Errorf("PathMetadataEntries length after append = %v, want 1", len(pm.PathMetadataEntries))
	}
}

// TestPackageMetadata_ContainsSpecialFiles tests that PackageMetadata contains SpecialFiles map.
func TestPackageMetadata_ContainsSpecialFiles(t *testing.T) {
	pm := NewPackageMetadata()
	if pm.SpecialFiles == nil {
		t.Fatal("SpecialFiles is nil")
	}
	if len(pm.SpecialFiles) != 0 {
		t.Errorf("SpecialFiles length = %v, want 0", len(pm.SpecialFiles))
	}
	// Test that we can add to SpecialFiles
	pm.SpecialFiles[65000] = nil
	if len(pm.SpecialFiles) != 1 {
		t.Errorf("SpecialFiles length after add = %v, want 1", len(pm.SpecialFiles))
	}
}
