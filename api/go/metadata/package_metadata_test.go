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

func runContainsSliceTest(t *testing.T, fieldName string, checkNil func(*PackageMetadata) bool, getLen func(*PackageMetadata) int, appendNil func(*PackageMetadata)) {
	t.Helper()
	pm := NewPackageMetadata()
	if checkNil(pm) {
		t.Fatalf("%s is nil", fieldName)
	}
	if getLen(pm) != 0 {
		t.Errorf("%s length = %v, want 0", fieldName, getLen(pm))
	}
	appendNil(pm)
	if getLen(pm) != 1 {
		t.Errorf("%s length after append = %v, want 1", fieldName, getLen(pm))
	}
}

type containsSliceCase struct {
	name      string
	fieldName string
	checkNil  func(*PackageMetadata) bool
	getLen    func(*PackageMetadata) int
	appendNil func(*PackageMetadata)
}

func makeContainsSliceCase(name, fieldName string, checkNil func(*PackageMetadata) bool, getLen func(*PackageMetadata) int, appendNil func(*PackageMetadata)) containsSliceCase {
	return containsSliceCase{name, fieldName, checkNil, getLen, appendNil}
}

//nolint:dupl // two slice fields require separate case builders with different accessors
func containsSliceCasesForTest() []containsSliceCase {
	fileEntriesCase := makeContainsSliceCase("FileEntries", "FileEntries",
		func(pm *PackageMetadata) bool { return pm.FileEntries == nil },
		func(pm *PackageMetadata) int { return len(pm.FileEntries) },
		func(pm *PackageMetadata) { pm.FileEntries = append(pm.FileEntries, nil) })
	pathMetaCase := makeContainsSliceCase("PathMetadataEntries", "PathMetadataEntries",
		func(pm *PackageMetadata) bool { return pm.PathMetadataEntries == nil },
		func(pm *PackageMetadata) int { return len(pm.PathMetadataEntries) },
		func(pm *PackageMetadata) { pm.PathMetadataEntries = append(pm.PathMetadataEntries, nil) })
	return []containsSliceCase{fileEntriesCase, pathMetaCase}
}

// TestPackageMetadata_ContainsSlices tests that PackageMetadata contains FileEntries and PathMetadataEntries slices.
func TestPackageMetadata_ContainsSlices(t *testing.T) {
	for _, tt := range containsSliceCasesForTest() {
		t.Run(tt.name, func(t *testing.T) {
			runContainsSliceTest(t, tt.fieldName, tt.checkNil, tt.getLen, tt.appendNil)
		})
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
