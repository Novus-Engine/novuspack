// Package metadata provides metadata domain structures for the NovusPack implementation.
//
// This file contains unit tests for PackageInfo structure.
package metadata

import (
	"testing"
	"time"

	"github.com/novus-engine/novuspack/api/go/signatures"
)

// TestPackageInfo tests the PackageInfo struct.
func TestPackageInfo(t *testing.T) {
	now := time.Now()
	signatureInfo := signatures.SignatureInfo{
		Index: 0,
		Type:  1,
		Valid: true,
	}

	info := PackageInfo{
		FileCount:               10,
		FilesUncompressedSize:   1024,
		FilesCompressedSize:     512,
		VendorID:                1,
		AppID:                   100,
		HasComment:              true,
		Comment:                 "Test comment",
		HasSignatures:           true,
		SignatureCount:          1,
		Signatures:              []signatures.SignatureInfo{signatureInfo},
		SecurityLevel:           SecurityLevelHigh,
		IsImmutable:             true,
		Created:                 now,
		Modified:                now,
		HasMetadataFiles:        true,
		HasEncryptedData:        false,
		HasCompressedData:       true,
		IsMetadataOnly:          false,
		PackageCompression:      1,
		IsPackageCompressed:     true,
		PackageOriginalSize:     2048,
		PackageCompressedSize:   1024,
		PackageCompressionRatio: 0.5,
	}

	if info.FileCount != 10 {
		t.Errorf("PackageInfo.FileCount = %v, want 10", info.FileCount)
	}
	if info.FilesUncompressedSize != 1024 {
		t.Errorf("PackageInfo.FilesUncompressedSize = %v, want 1024", info.FilesUncompressedSize)
	}
	if info.FilesCompressedSize != 512 {
		t.Errorf("PackageInfo.FilesCompressedSize = %v, want 512", info.FilesCompressedSize)
	}
	if info.VendorID != 1 {
		t.Errorf("PackageInfo.VendorID = %v, want 1", info.VendorID)
	}
	if info.AppID != 100 {
		t.Errorf("PackageInfo.AppID = %v, want 100", info.AppID)
	}
	if !info.HasComment {
		t.Errorf("PackageInfo.HasComment = %v, want true", info.HasComment)
	}
	if info.Comment != "Test comment" {
		t.Errorf("PackageInfo.Comment = %v, want Test comment", info.Comment)
	}
	if !info.HasSignatures {
		t.Errorf("PackageInfo.HasSignatures = %v, want true", info.HasSignatures)
	}
	if info.SignatureCount != 1 {
		t.Errorf("PackageInfo.SignatureCount = %v, want 1", info.SignatureCount)
	}
	if len(info.Signatures) != 1 {
		t.Errorf("PackageInfo.Signatures length = %v, want 1", len(info.Signatures))
	}
	if info.SecurityLevel != SecurityLevelHigh {
		t.Errorf("PackageInfo.SecurityLevel = %v, want SecurityLevelHigh", info.SecurityLevel)
	}
	if !info.IsImmutable {
		t.Errorf("PackageInfo.IsImmutable = %v, want true", info.IsImmutable)
	}
	if !info.Created.Equal(now) {
		t.Errorf("PackageInfo.Created = %v, want %v", info.Created, now)
	}
	if !info.Modified.Equal(now) {
		t.Errorf("PackageInfo.Modified = %v, want %v", info.Modified, now)
	}
	if !info.HasMetadataFiles {
		t.Errorf("PackageInfo.HasMetadataFiles = %v, want true", info.HasMetadataFiles)
	}
	if info.HasEncryptedData {
		t.Errorf("PackageInfo.HasEncryptedData = %v, want false", info.HasEncryptedData)
	}
	if !info.HasCompressedData {
		t.Errorf("PackageInfo.HasCompressedData = %v, want true", info.HasCompressedData)
	}
	if info.IsMetadataOnly {
		t.Errorf("PackageInfo.IsMetadataOnly = %v, want false", info.IsMetadataOnly)
	}
	if info.PackageCompression != 1 {
		t.Errorf("PackageInfo.PackageCompression = %v, want 1", info.PackageCompression)
	}
	if !info.IsPackageCompressed {
		t.Errorf("PackageInfo.IsPackageCompressed = %v, want true", info.IsPackageCompressed)
	}
	if info.PackageOriginalSize != 2048 {
		t.Errorf("PackageInfo.PackageOriginalSize = %v, want 2048", info.PackageOriginalSize)
	}
	if info.PackageCompressedSize != 1024 {
		t.Errorf("PackageInfo.PackageCompressedSize = %v, want 1024", info.PackageCompressedSize)
	}
	if info.PackageCompressionRatio != 0.5 {
		t.Errorf("PackageInfo.PackageCompressionRatio = %v, want 0.5", info.PackageCompressionRatio)
	}
}

// TestNewPackageInfo tests the NewPackageInfo function.
func TestNewPackageInfo(t *testing.T) {
	info := NewPackageInfo()

	if info == nil {
		t.Fatal("NewPackageInfo() returned nil")
	}

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so info is not nil after check
	if info.FileCount != 0 {
		t.Errorf("NewPackageInfo() FileCount = %d, want 0", info.FileCount)
	}

	if info.FilesUncompressedSize != 0 {
		t.Errorf("NewPackageInfo() FilesUncompressedSize = %d, want 0", info.FilesUncompressedSize)
	}

	if info.FilesCompressedSize != 0 {
		t.Errorf("NewPackageInfo() FilesCompressedSize = %d, want 0", info.FilesCompressedSize)
	}

	if info.VendorID != 0 {
		t.Errorf("NewPackageInfo() VendorID = %d, want 0", info.VendorID)
	}

	if info.AppID != 0 {
		t.Errorf("NewPackageInfo() AppID = %d, want 0", info.AppID)
	}

	if info.HasComment {
		t.Error("NewPackageInfo() HasComment = true, want false")
	}

	if info.Comment != "" {
		t.Errorf("NewPackageInfo() Comment = %q, want empty", info.Comment)
	}

	if info.HasSignatures {
		t.Error("NewPackageInfo() HasSignatures = true, want false")
	}

	if info.SignatureCount != 0 {
		t.Errorf("NewPackageInfo() SignatureCount = %d, want 0", info.SignatureCount)
	}

	if len(info.Signatures) != 0 {
		t.Errorf("NewPackageInfo() Signatures length = %d, want 0", len(info.Signatures))
	}

	if info.SecurityLevel != SecurityLevelNone {
		t.Errorf("NewPackageInfo() SecurityLevel = %v, want SecurityLevelNone", info.SecurityLevel)
	}

	if info.IsImmutable {
		t.Error("NewPackageInfo() IsImmutable = true, want false")
	}

	if !info.Created.IsZero() {
		t.Error("NewPackageInfo() Created should be zero time")
	}

	if !info.Modified.IsZero() {
		t.Error("NewPackageInfo() Modified should be zero time")
	}

	if info.HasMetadataFiles {
		t.Error("NewPackageInfo() HasMetadataFiles = true, want false")
	}

	if info.HasEncryptedData {
		t.Error("NewPackageInfo() HasEncryptedData = true, want false")
	}

	if info.HasCompressedData {
		t.Error("NewPackageInfo() HasCompressedData = true, want false")
	}

	if info.IsMetadataOnly {
		t.Error("NewPackageInfo() IsMetadataOnly = true, want false")
	}

	if info.PackageCompression != 0 {
		t.Errorf("NewPackageInfo() PackageCompression = %d, want 0", info.PackageCompression)
	}

	if info.IsPackageCompressed {
		t.Error("NewPackageInfo() IsPackageCompressed = true, want false")
	}

	if info.PackageOriginalSize != 0 {
		t.Errorf("NewPackageInfo() PackageOriginalSize = %d, want 0", info.PackageOriginalSize)
	}

	if info.PackageCompressedSize != 0 {
		t.Errorf("NewPackageInfo() PackageCompressedSize = %d, want 0", info.PackageCompressedSize)
	}

	if info.PackageCompressionRatio != 0.0 {
		t.Errorf("NewPackageInfo() PackageCompressionRatio = %v, want 0.0", info.PackageCompressionRatio)
	}
}

// TestPackageInfo_ZeroValue tests the zero value of PackageInfo.
func TestPackageInfo_ZeroValue(t *testing.T) {
	info := PackageInfo{}

	if info.FileCount != 0 {
		t.Errorf("PackageInfo zero value FileCount = %v, want 0", info.FileCount)
	}
	if info.FilesUncompressedSize != 0 {
		t.Errorf("PackageInfo zero value FilesUncompressedSize = %v, want 0", info.FilesUncompressedSize)
	}
	if info.FilesCompressedSize != 0 {
		t.Errorf("PackageInfo zero value FilesCompressedSize = %v, want 0", info.FilesCompressedSize)
	}
	if info.VendorID != 0 {
		t.Errorf("PackageInfo zero value VendorID = %v, want 0", info.VendorID)
	}
	if info.AppID != 0 {
		t.Errorf("PackageInfo zero value AppID = %v, want 0", info.AppID)
	}
	if info.HasComment {
		t.Errorf("PackageInfo zero value HasComment = %v, want false", info.HasComment)
	}
	if info.Comment != "" {
		t.Errorf("PackageInfo zero value Comment = %v, want empty string", info.Comment)
	}
	if info.HasSignatures {
		t.Errorf("PackageInfo zero value HasSignatures = %v, want false", info.HasSignatures)
	}
	if info.SignatureCount != 0 {
		t.Errorf("PackageInfo zero value SignatureCount = %v, want 0", info.SignatureCount)
	}
	if info.Signatures != nil {
		t.Errorf("PackageInfo zero value Signatures = %v, want nil", info.Signatures)
	}
	if info.SecurityLevel != SecurityLevelNone {
		t.Errorf("PackageInfo zero value SecurityLevel = %v, want SecurityLevelNone", info.SecurityLevel)
	}
	if info.IsImmutable {
		t.Errorf("PackageInfo zero value IsImmutable = %v, want false", info.IsImmutable)
	}
	if !info.Created.IsZero() {
		t.Errorf("PackageInfo zero value Created = %v, want zero time", info.Created)
	}
	if !info.Modified.IsZero() {
		t.Errorf("PackageInfo zero value Modified = %v, want zero time", info.Modified)
	}
	if info.HasMetadataFiles {
		t.Errorf("PackageInfo zero value HasMetadataFiles = %v, want false", info.HasMetadataFiles)
	}
	if info.HasEncryptedData {
		t.Errorf("PackageInfo zero value HasEncryptedData = %v, want false", info.HasEncryptedData)
	}
	if info.HasCompressedData {
		t.Errorf("PackageInfo zero value HasCompressedData = %v, want false", info.HasCompressedData)
	}
	if info.IsMetadataOnly {
		t.Errorf("PackageInfo zero value IsMetadataOnly = %v, want false", info.IsMetadataOnly)
	}
	if info.PackageCompression != 0 {
		t.Errorf("PackageInfo zero value PackageCompression = %v, want 0", info.PackageCompression)
	}
	if info.IsPackageCompressed {
		t.Errorf("PackageInfo zero value IsPackageCompressed = %v, want false", info.IsPackageCompressed)
	}
	if info.PackageOriginalSize != 0 {
		t.Errorf("PackageInfo zero value PackageOriginalSize = %v, want 0", info.PackageOriginalSize)
	}
	if info.PackageCompressedSize != 0 {
		t.Errorf("PackageInfo zero value PackageCompressedSize = %v, want 0", info.PackageCompressedSize)
	}
	if info.PackageCompressionRatio != 0.0 {
		t.Errorf("PackageInfo zero value PackageCompressionRatio = %v, want 0.0", info.PackageCompressionRatio)
	}
}
