package fileformat

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// packageHeaderMinimal returns a PackageHeader with minimal/default field values (Flags 0, ArchivePartInfo Part 1 of 1).
func packageHeaderMinimal() PackageHeader {
	return PackageHeader{
		Magic:              NVPKMagic,
		FormatVersion:      FormatVersion,
		Flags:              0,
		PackageDataVersion: 1,
		MetadataVersion:    1,
		PackageCRC:         0,
		CreatedTime:        0,
		ModifiedTime:       0,
		LocaleID:           0,
		Reserved:           0,
		AppID:              0,
		VendorID:           0,
		CreatorID:          0,
		IndexStart:         0,
		IndexSize:          0,
		ArchiveChainID:     0,
		ArchivePartInfo:    0x00010001, // Part 1 of 1
		CommentSize:        0,
		CommentStart:       0,
		SignatureOffset:    0,
	}
}

// packageHeaderMinimalZeroPart returns a PackageHeader with minimal fields and ArchivePartInfo 0.
func packageHeaderMinimalZeroPart() PackageHeader {
	h := packageHeaderMinimal()
	h.ArchivePartInfo = 0
	return h
}

// packageHeaderForSerialization returns a full PackageHeader used in serialization/ReadFrom tests.
func packageHeaderForSerialization() PackageHeader {
	return PackageHeader{
		Magic:              NVPKMagic,
		FormatVersion:      FormatVersion,
		Flags:              FlagHasSignatures | FlagHasCompressedFiles,
		PackageDataVersion: 1,
		MetadataVersion:    1,
		PackageCRC:         0x12345678,
		CreatedTime:        1638360000000000000,
		ModifiedTime:       1638360000000000000,
		LocaleID:           0x0409, // en-US
		Reserved:           0,
		AppID:              730, // CS:GO
		VendorID:           VendorIDSteam,
		CreatorID:          0,
		IndexStart:         4096,
		IndexSize:          1024,
		ArchiveChainID:     0,
		ArchivePartInfo:    0x00010001,
		CommentSize:        0,
		CommentStart:       0,
		SignatureOffset:    0,
	}
}

// packageHeaderWithAllFieldsSet returns a PackageHeader with all fields set (for WriteTo/ReadFrom tests).
func packageHeaderWithAllFieldsSet() PackageHeader {
	return packageHeaderFull(42, 17, 0xDEADBEEF, 0x00020003, 100, 0x0411)
}

// packageHeaderRoundTripFull returns a full PackageHeader for round-trip tests (different metadata values).
func packageHeaderRoundTripFull() PackageHeader {
	return packageHeaderFull(100, 50, 0xABCDEF00, 0x0005000A, 200, 0x0409)
}

func packageHeaderFull(pkgDataVer, metaVer, crc, partInfo, commentSize, localeID uint32) PackageHeader {
	return PackageHeader{
		Magic:              NVPKMagic,
		FormatVersion:      FormatVersion,
		Flags:              0x01FF,
		PackageDataVersion: pkgDataVer,
		MetadataVersion:    metaVer,
		PackageCRC:         crc,
		CreatedTime:        1638360000000000000,
		ModifiedTime:       1638361000000000000,
		LocaleID:           localeID,
		Reserved:           0,
		AppID:              730,
		VendorID:           VendorIDSteam,
		CreatorID:          0,
		IndexStart:         8192,
		IndexSize:          2048,
		ArchiveChainID:     0x123456789ABCDEF0,
		ArchivePartInfo:    partInfo,
		CommentSize:        commentSize,
		CommentStart:       6144,
		SignatureOffset:    10240,
	}
}

// headerGetterTestCase is a single test case for a getter on PackageHeader.
type headerGetterTestCase[T comparable] struct {
	name   string
	header PackageHeader
	want   T
}

// runHeaderGetterTests runs table-driven tests for a header getter; T must be comparable.
func runHeaderGetterTests[T comparable](t *testing.T, tests []headerGetterTestCase[T], getter func(PackageHeader) T, format string) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getter(tt.header)
			if got != tt.want {
				t.Errorf(format, got, tt.want)
			}
		})
	}
}

// TestPackageHeaderSize verifies the PackageHeader struct is exactly 112 bytes
// Specification: package_file_format.md: 2.1 Header Structure
func TestPackageHeaderSize(t *testing.T) {
	var header PackageHeader
	size := binary.Size(header)

	if size != PackageHeaderSize {
		t.Errorf("PackageHeader size = %d bytes, want %d bytes", size, PackageHeaderSize)
	}
}

// TestPackageHeaderFieldTypes verifies all fields have correct types
// Specification: package_file_format.md: 2.1 Header Structure
func TestPackageHeaderFieldTypes(t *testing.T) {
	header := packageHeaderMinimal()

	// Verify Magic is uint32
	if header.Magic != NVPKMagic {
		t.Errorf("Magic = 0x%X, want 0x%X", header.Magic, NVPKMagic)
	}

	// Verify FormatVersion is uint32
	if header.FormatVersion != FormatVersion {
		t.Errorf("FormatVersion = %d, want %d", header.FormatVersion, FormatVersion)
	}

	// Verify ArchivePartInfo default value
	if header.ArchivePartInfo != 0x00010001 {
		t.Errorf("ArchivePartInfo = 0x%X, want 0x00010001", header.ArchivePartInfo)
	}
}

// TestPackageHeaderSerialization verifies binary serialization/deserialization
// Specification: package_file_format.md: 2.1 Header Structure
func TestPackageHeaderSerialization(t *testing.T) {
	original := packageHeaderForSerialization()

	// Serialize to bytes
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &original)
	if err != nil {
		t.Fatalf("Failed to serialize header: %v", err)
	}

	// Verify size
	if buf.Len() != PackageHeaderSize {
		t.Errorf("Serialized header size = %d bytes, want %d bytes", buf.Len(), PackageHeaderSize)
	}

	// Deserialize from bytes
	var deserialized PackageHeader
	err = binary.Read(buf, binary.LittleEndian, &deserialized)
	if err != nil {
		t.Fatalf("Failed to deserialize header: %v", err)
	}

	// Verify all fields match
	if deserialized != original {
		t.Errorf("Deserialized header does not match original")
		t.Logf("Original: %+v", original)
		t.Logf("Deserialized: %+v", deserialized)
	}
}

// TestPackageHeaderMagicValidation verifies magic number validation
// Specification: package_file_format.md: 2.1 Header Structure
func TestPackageHeaderMagicValidation(t *testing.T) {
	tests := []struct {
		name    string
		magic   uint32
		version uint32
		wantErr bool
	}{
		{"Valid magic", NVPKMagic, FormatVersion, false},
		{"Invalid magic", 0x12345678, FormatVersion, true},
		{"Zero magic", 0x00000000, FormatVersion, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := PackageHeader{
				Magic:         tt.magic,
				FormatVersion: tt.version,
			}
			err := header.validate()

			if tt.wantErr && err == nil {
				t.Error("Validate() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}

// TestPackageHeaderReservedFieldValidation verifies reserved field must be zero
// Specification: package_file_format.md: 2.1 Header Structure
func TestPackageHeaderReservedFieldValidation(t *testing.T) {
	header := PackageHeader{
		Magic:         NVPKMagic,
		FormatVersion: FormatVersion,
		Reserved:      1, // Non-zero reserved field
	}

	err := header.validate()
	if err == nil {
		t.Error("Validate() expected error for non-zero reserved field, got nil")
	}
}

// TestPackageHeaderFormatVersionValidation verifies format version validation
// Specification: package_file_format.md: 2.1 Header Structure
func TestPackageHeaderFormatVersionValidation(t *testing.T) {
	tests := []struct {
		name          string
		formatVersion uint32
		wantErr       bool
	}{
		{"Valid format version", FormatVersion, false},
		{"Invalid format version", 0, true},
		{"Invalid format version", 2, true},
		{"Invalid format version", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := PackageHeader{
				Magic:         NVPKMagic,
				FormatVersion: tt.formatVersion,
			}
			err := header.validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestArchivePartInfoEncoding verifies ArchivePartInfo packing/unpacking
// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
func TestArchivePartInfoEncoding(t *testing.T) {
	tests := []struct {
		name      string
		partInfo  uint32
		wantPart  uint16
		wantTotal uint16
	}{
		{"Single archive", 0x00010001, 1, 1},
		{"Part 2 of 3", 0x00020003, 2, 3},
		{"Part 0 of 0", 0x00000000, 0, 0},
		{"Part 65535 of 65535", 0xFFFFFFFF, 65535, 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := PackageHeader{ArchivePartInfo: tt.partInfo}

			part := uint16(header.ArchivePartInfo >> 16)
			total := uint16(header.ArchivePartInfo & 0xFFFF)

			if part != tt.wantPart {
				t.Errorf("Part number = %d, want %d", part, tt.wantPart)
			}
			if total != tt.wantTotal {
				t.Errorf("Total parts = %d, want %d", total, tt.wantTotal)
			}
		})
	}
}

// TestPackageHeaderFlagsEncoding verifies Flags field encoding/decoding
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderFlagsEncoding(t *testing.T) {
	tests := []struct {
		name            string
		flags           uint32
		wantCompression uint8
		wantFeatures    uint8
	}{
		{"No compression, no features", 0x0000, 0, 0x00},
		{"Zstd compression", 0x0100, 1, 0x00},
		{"LZ4 compression", 0x0200, 2, 0x00},
		{"LZMA compression", 0x0300, 3, 0x00},
		{"Has signatures", 0x0001, 0, 0x01},
		{"Has signatures and compressed", 0x0003, 0, 0x03},
		{"Zstd with all features", 0x01FF, 1, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := PackageHeader{Flags: tt.flags}

			compression := uint8((header.Flags & FlagsMaskCompressionType) >> FlagsShiftCompressionType)
			features := uint8(header.Flags & FlagsMaskFeatures)

			if compression != tt.wantCompression {
				t.Errorf("Compression type = %d, want %d", compression, tt.wantCompression)
			}
			if features != tt.wantFeatures {
				t.Errorf("Features = 0x%02X, want 0x%02X", features, tt.wantFeatures)
			}
		})
	}
}

// TestPackageHeaderVersionFields verifies version field behavior
// Specification: package_file_format.md: 2.2 Package Version Fields Specification
func TestPackageHeaderVersionFields(t *testing.T) {
	header := PackageHeader{
		Magic:              NVPKMagic,
		FormatVersion:      FormatVersion,
		PackageDataVersion: 1,
		MetadataVersion:    1,
	}

	// Verify initial values
	if header.PackageDataVersion != 1 {
		t.Errorf("PackageDataVersion = %d, want 1", header.PackageDataVersion)
	}
	if header.MetadataVersion != 1 {
		t.Errorf("MetadataVersion = %d, want 1", header.MetadataVersion)
	}

	// Simulate version increment
	header.PackageDataVersion++
	if header.PackageDataVersion != 2 {
		t.Errorf("PackageDataVersion after increment = %d, want 2", header.PackageDataVersion)
	}

	header.MetadataVersion++
	if header.MetadataVersion != 2 {
		t.Errorf("MetadataVersion after increment = %d, want 2", header.MetadataVersion)
	}
}

// TestPackageHeaderGetCompressionType verifies GetCompressionType extraction
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderGetCompressionType(t *testing.T) {
	tests := []headerGetterTestCase[uint8]{
		{"No compression", PackageHeader{Flags: 0x0000}, CompressionNone},
		{"Zstd compression", PackageHeader{Flags: 0x0100}, CompressionZstd},
		{"LZ4 compression", PackageHeader{Flags: 0x0200}, CompressionLZ4},
		{"LZMA compression", PackageHeader{Flags: 0x0300}, CompressionLZMA},
		{"Zstd with features", PackageHeader{Flags: 0x01FF}, CompressionZstd},
		{"LZ4 with features", PackageHeader{Flags: 0x02FF}, CompressionLZ4},
	}
	runHeaderGetterTests(t, tests, func(h PackageHeader) uint8 { return h.getCompressionType() }, "getCompressionType() = %d, want %d")
}

// TestPackageHeaderSetCompressionType verifies SetCompressionType preserves features
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderSetCompressionType(t *testing.T) {
	tests := []struct {
		name            string
		initialFlags    uint32
		compressionType uint8
		wantCompression uint8
		wantFeatures    uint8
	}{
		{"Set Zstd on empty flags", 0x0000, CompressionZstd, CompressionZstd, 0x00},
		{"Set LZ4 on empty flags", 0x0000, CompressionLZ4, CompressionLZ4, 0x00},
		{"Set Zstd preserves features", 0x00FF, CompressionZstd, CompressionZstd, 0xFF},
		{"Set LZ4 preserves features", 0x00FF, CompressionLZ4, CompressionLZ4, 0xFF},
		{"Change from Zstd to LZ4", 0x01FF, CompressionLZ4, CompressionLZ4, 0xFF},
		{"Change from LZ4 to LZMA", 0x02FF, CompressionLZMA, CompressionLZMA, 0xFF},
		{"Set None clears compression", 0x01FF, CompressionNone, CompressionNone, 0xFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := PackageHeader{Flags: tt.initialFlags}
			header.setCompressionType(tt.compressionType)

			gotCompression := header.getCompressionType()
			gotFeatures := header.getFeatures()

			if gotCompression != tt.wantCompression {
				t.Errorf("Compression type = %d, want %d", gotCompression, tt.wantCompression)
			}
			if gotFeatures != tt.wantFeatures {
				t.Errorf("Features = 0x%02X, want 0x%02X", gotFeatures, tt.wantFeatures)
			}
		})
	}
}

// TestPackageHeaderGetFeatures verifies GetFeatures extraction
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderGetFeatures(t *testing.T) {
	tests := []headerGetterTestCase[uint8]{
		{"No features", PackageHeader{Flags: 0x0000}, 0x00},
		{"Has signatures", PackageHeader{Flags: FlagHasSignatures}, 0x01},
		{"Has compressed files", PackageHeader{Flags: FlagHasCompressedFiles}, 0x02},
		{"Has encrypted files", PackageHeader{Flags: FlagHasEncryptedFiles}, 0x04},
		{"All features", PackageHeader{Flags: 0x00FF}, 0xFF},
		{"Features with compression", PackageHeader{Flags: 0x01FF}, 0xFF},
	}
	runHeaderGetterTests(t, tests, func(h PackageHeader) uint8 { return h.getFeatures() }, "getFeatures() = 0x%02X, want 0x%02X")
}

// TestPackageHeaderHasFeature verifies HasFeature checking
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderHasFeature(t *testing.T) {
	header := PackageHeader{Flags: FlagHasSignatures | FlagHasCompressedFiles}

	tests := []struct {
		name string
		flag uint32
		want bool
	}{
		{"Has signatures", FlagHasSignatures, true},
		{"Has compressed files", FlagHasCompressedFiles, true},
		{"Has encrypted files", FlagHasEncryptedFiles, false},
		{"Has package comment", FlagHasPackageComment, false},
		{"Has per-file tags", FlagHasPerFileTags, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := header.hasFeature(tt.flag)
			if got != tt.want {
				t.Errorf("HasFeature(0x%08X) = %v, want %v", tt.flag, got, tt.want)
			}
		})
	}
}

// TestPackageHeaderSetFeature verifies SetFeature setting
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderSetFeature(t *testing.T) {
	header := PackageHeader{Flags: 0x0000}

	// Set first feature
	header.setFeature(FlagHasSignatures)
	if !header.hasFeature(FlagHasSignatures) {
		t.Error("SetFeature(FlagHasSignatures) did not set the flag")
	}

	// Set additional feature
	header.setFeature(FlagHasCompressedFiles)
	if !header.hasFeature(FlagHasSignatures) {
		t.Error("SetFeature(FlagHasCompressedFiles) cleared existing flag")
	}
	if !header.hasFeature(FlagHasCompressedFiles) {
		t.Error("SetFeature(FlagHasCompressedFiles) did not set the flag")
	}

	// Set multiple features
	header.setFeature(FlagHasEncryptedFiles | FlagHasPackageComment)
	if !header.hasFeature(FlagHasEncryptedFiles) {
		t.Error("SetFeature did not set FlagHasEncryptedFiles")
	}
	if !header.hasFeature(FlagHasPackageComment) {
		t.Error("SetFeature did not set FlagHasPackageComment")
	}
}

// TestPackageHeaderClearFeature verifies ClearFeature clearing
// Specification: package_file_format.md: 2.5 Package Features Flags
func TestPackageHeaderClearFeature(t *testing.T) {
	header := PackageHeader{Flags: FlagHasSignatures | FlagHasCompressedFiles | FlagHasEncryptedFiles}

	// Clear one feature
	header.clearFeature(FlagHasSignatures)
	if header.hasFeature(FlagHasSignatures) {
		t.Error("ClearFeature(FlagHasSignatures) did not clear the flag")
	}
	if !header.hasFeature(FlagHasCompressedFiles) {
		t.Error("ClearFeature(FlagHasSignatures) cleared wrong flag")
	}
	if !header.hasFeature(FlagHasEncryptedFiles) {
		t.Error("ClearFeature(FlagHasSignatures) cleared wrong flag")
	}

	// Clear multiple features
	header.clearFeature(FlagHasCompressedFiles | FlagHasEncryptedFiles)
	if header.hasFeature(FlagHasCompressedFiles) {
		t.Error("ClearFeature did not clear FlagHasCompressedFiles")
	}
	if header.hasFeature(FlagHasEncryptedFiles) {
		t.Error("ClearFeature did not clear FlagHasEncryptedFiles")
	}
}

// TestPackageHeaderArchivePartInfoGetters verifies GetArchivePart and GetArchiveTotal extraction.
// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
func TestPackageHeaderArchivePartInfoGetters(t *testing.T) {
	tests := []struct {
		name      string
		header    PackageHeader
		wantPart  uint16
		wantTotal uint16
	}{
		{"Part 1 of 1", PackageHeader{ArchivePartInfo: 0x00010001}, 1, 1},
		{"Part 2 of 3", PackageHeader{ArchivePartInfo: 0x00020003}, 2, 3},
		{"Part 0 of 0", PackageHeader{ArchivePartInfo: 0x00000000}, 0, 0},
		{"Part 65535 of 65535", PackageHeader{ArchivePartInfo: 0xFFFFFFFF}, 65535, 65535},
		{"Part 10 of 20", PackageHeader{ArchivePartInfo: 0x000A0014}, 10, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.header.getArchivePart(); got != tt.wantPart {
				t.Errorf("getArchivePart() = %d, want %d", got, tt.wantPart)
			}
			if got := tt.header.getArchiveTotal(); got != tt.wantTotal {
				t.Errorf("getArchiveTotal() = %d, want %d", got, tt.wantTotal)
			}
		})
	}
}

// TestPackageHeaderSetArchivePartInfo verifies SetArchivePartInfo setting
// Specification: package_file_format.md: 2.6 ArchivePartInfo Field Specification
func TestPackageHeaderSetArchivePartInfo(t *testing.T) {
	tests := []struct {
		name      string
		part      uint16
		total     uint16
		wantPart  uint16
		wantTotal uint16
	}{
		{"Part 1 of 1", 1, 1, 1, 1},
		{"Part 2 of 3", 2, 3, 2, 3},
		{"Part 0 of 0", 0, 0, 0, 0},
		{"Part 65535 of 65535", 65535, 65535, 65535, 65535},
		{"Part 10 of 20", 10, 20, 10, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := PackageHeader{}
			header.setArchivePartInfo(tt.part, tt.total)

			gotPart := header.getArchivePart()
			gotTotal := header.getArchiveTotal()

			if gotPart != tt.wantPart {
				t.Errorf("Part = %d, want %d", gotPart, tt.wantPart)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("Total = %d, want %d", gotTotal, tt.wantTotal)
			}
		})
	}
}

// isSignedGetterTestCases returns test cases for IsSigned().
func isSignedGetterTestCases() []headerGetterTestCase[bool] {
	return []headerGetterTestCase[bool]{
		{"Not signed", PackageHeader{SignatureOffset: 0}, false},
		{"Signed", PackageHeader{SignatureOffset: 4096}, true},
		{"Signed at offset 1", PackageHeader{SignatureOffset: 1}, true},
		{"Signed at large offset", PackageHeader{SignatureOffset: 0xFFFFFFFFFFFFFFFF}, true},
	}
}

// hasCommentGetterTestCases returns test cases for HasComment().
func hasCommentGetterTestCases() []headerGetterTestCase[bool] {
	return []headerGetterTestCase[bool]{
		{"No comment", PackageHeader{CommentSize: 0}, false},
		{"Has comment", PackageHeader{CommentSize: 100}, true},
		{"Has comment size 1", PackageHeader{CommentSize: 1}, true},
		{"Has large comment", PackageHeader{CommentSize: 0xFFFFFFFF}, true},
	}
}

// TestPackageHeaderIsSignedAndHasComment verifies IsSigned and HasComment checking.
// Specification: package_file_format.md: 2.9 2.9 Signed Package File Immutability and Incremental Signatures
// Specification: package_file_format.md: 7.1 7.1 Package Comment Format Specification
func TestPackageHeaderIsSignedAndHasComment(t *testing.T) {
	t.Run("IsSigned", func(t *testing.T) {
		runHeaderGetterTests(t, isSignedGetterTestCases(), func(h PackageHeader) bool { return h.isSigned() }, "isSigned() = %v, want %v")
	})
	t.Run("HasComment", func(t *testing.T) {
		runHeaderGetterTests(t, hasCommentGetterTestCases(), func(h PackageHeader) bool { return h.hasComment() }, "hasComment() = %v, want %v")
	})
}

// TestNewPackageHeader verifies NewPackageHeader initializes correctly
// Specification: package_file_format.md: 2.8.1 Initial Package Creation
func TestNewPackageHeader(t *testing.T) {
	header := NewPackageHeader()

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so header is not nil after check
	if header == nil {
		t.Fatal("NewPackageHeader() returned nil")
	}

	// Verify default values per spec Section 2.8.1
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so header is not nil after check
	if header.Magic != NVPKMagic {
		t.Errorf("Magic = 0x%08X, want 0x%08X", header.Magic, NVPKMagic)
	}
	if header.FormatVersion != FormatVersion {
		t.Errorf("FormatVersion = %d, want %d", header.FormatVersion, FormatVersion)
	}
	if header.PackageDataVersion != 1 {
		t.Errorf("PackageDataVersion = %d, want 1", header.PackageDataVersion)
	}
	if header.MetadataVersion != 1 {
		t.Errorf("MetadataVersion = %d, want 1", header.MetadataVersion)
	}
	if header.Reserved != 0 {
		t.Errorf("Reserved = %d, want 0", header.Reserved)
	}
	if header.ArchivePartInfo != 0x00010001 {
		t.Errorf("ArchivePartInfo = 0x%08X, want 0x00010001", header.ArchivePartInfo)
	}
	if header.CreatorID != 0 {
		t.Errorf("CreatorID = %d, want 0", header.CreatorID)
	}
	if header.CommentSize != 0 {
		t.Errorf("CommentSize = %d, want 0", header.CommentSize)
	}
	if header.CommentStart != 0 {
		t.Errorf("CommentStart = %d, want 0", header.CommentStart)
	}
	if header.SignatureOffset != 0 {
		t.Errorf("SignatureOffset = %d, want 0", header.SignatureOffset)
	}

	// Verify it passes validation
	if err := header.validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

// TestPackageHeaderReadFrom verifies ReadFrom deserialization
// Specification: package_file_format.md: 2.1 Header Structure
//
//nolint:gocognit // table-driven test
func TestPackageHeaderReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		header  PackageHeader
		wantErr bool
	}{
		{"Valid header", packageHeaderForSerialization(), false},
		{"Header with all fields set", packageHeaderWithAllFieldsSet(), false},
		{"Header with minimal fields", packageHeaderMinimalZeroPart(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Serialize original header
			buf := new(bytes.Buffer)
			err := binary.Write(buf, binary.LittleEndian, &tt.header)
			if err != nil {
				t.Fatalf("Failed to serialize header: %v", err)
			}

			// Verify size
			if buf.Len() != PackageHeaderSize {
				t.Errorf("Serialized header size = %d bytes, want %d bytes", buf.Len(), PackageHeaderSize)
			}

			// Deserialize using ReadFrom
			var header PackageHeader
			n, err := header.readFrom(buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if n != PackageHeaderSize {
					t.Errorf("ReadFrom() read %d bytes, want %d", n, PackageHeaderSize)
				}

				// Verify all fields match
				if header != tt.header {
					t.Errorf("ReadFrom() header mismatch")
					t.Logf("Original: %+v", tt.header)
					t.Logf("Read: %+v", header)
				}

				// Verify validation passes
				if err := header.validate(); err != nil {
					t.Errorf("ReadFrom() header validation failed: %v", err)
				}
			}
		})
	}
}

// TestPackageHeaderReadFromInvalidMagic verifies ReadFrom rejects invalid magic
func TestPackageHeaderReadFromInvalidMagic(t *testing.T) {
	// Create header with invalid magic
	invalidHeader := PackageHeader{
		Magic:         0x12345678, // Invalid magic
		FormatVersion: FormatVersion,
		Reserved:      0,
	}

	// Serialize
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &invalidHeader)
	if err != nil {
		t.Fatalf("Failed to serialize header: %v", err)
	}

	// Try to read
	var header PackageHeader
	_, err = header.readFrom(buf)
	if err == nil {
		t.Error("ReadFrom() expected error for invalid magic, got nil")
	} else if !strings.Contains(err.Error(), "magic") {
		t.Errorf("ReadFrom() error = %q, want error containing 'magic'", err.Error())
	}
}

// TestPackageHeaderReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestPackageHeaderReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Partial data", make([]byte, 50)},
		{"Almost complete", make([]byte, 111)},
		{"Partial header with invalid magic", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(0x12345678)) // Invalid magic
			// Rest incomplete
			return buf.Bytes()[:50]
		}()},
		{"Complete header with invalid magic", func() []byte {
			header := PackageHeader{
				Magic:              0x12345678, // Invalid magic
				FormatVersion:      FormatVersion,
				PackageDataVersion: 1,
				MetadataVersion:    1,
			}
			var buf bytes.Buffer
			_ = binary.Write(&buf, binary.LittleEndian, &header)
			return buf.Bytes()
		}()},
		{"Header with various partial reads", func() []byte {
			// Test different partial read scenarios
			return make([]byte, 30) // Very partial
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var header PackageHeader
			r := bytes.NewReader(tt.data)
			_, err := header.readFrom(r)

			// Check if this is a valid case (complete header with invalid magic)
			isInvalidMagicCase := strings.Contains(tt.name, "Complete header with invalid magic")
			if isInvalidMagicCase {
				if err == nil {
					t.Errorf("ReadFrom() expected error for invalid magic, got nil")
				} else if !strings.Contains(err.Error(), "magic") {
					t.Errorf("ReadFrom() error = %q, want error containing 'magic'", err.Error())
				}
			} else {
				if err == nil {
					t.Errorf("ReadFrom() expected error for incomplete data, got nil")
				}
			}
		})
	}
}

// TestPackageHeaderReadFromNonEOFError verifies ReadFrom handles non-EOF errors
func TestPackageHeaderReadFromNonEOFError(t *testing.T) {
	var header PackageHeader
	r := testhelpers.NewErrorReader()
	_, err := header.readFrom(r)

	switch {
	case err == nil:
		t.Error("ReadFrom() expected error for error reader, got nil")
	case strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "incomplete"):
		t.Errorf("ReadFrom() error = %q, want non-EOF error", err.Error())
	case !strings.Contains(err.Error(), "failed to read header"):
		t.Errorf("ReadFrom() error = %q, want error containing 'failed to read header'", err.Error())
	}
}

// TestPackageHeaderWriteTo verifies WriteTo serialization
// Specification: package_file_format.md: 2.1 Header Structure
//
//nolint:gocognit // table-driven test
func TestPackageHeaderWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		header  PackageHeader
		wantErr bool
	}{
		{
			"Valid header",
			PackageHeader{
				Magic:              NVPKMagic,
				FormatVersion:      FormatVersion,
				Flags:              FlagHasSignatures,
				PackageDataVersion: 1,
				MetadataVersion:    1,
				Reserved:           0,
				ArchivePartInfo:    0x00010001,
			},
			false,
		},
		{"Header with all fields", packageHeaderWithAllFieldsSet(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.header.writeTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if n != PackageHeaderSize {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, PackageHeaderSize)
				}

				if buf.Len() != PackageHeaderSize {
					t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), PackageHeaderSize)
				}

				// Verify we can read it back
				var header PackageHeader
				_, readErr := header.readFrom(&buf)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if header != tt.header {
					t.Errorf("Read back header does not match original")
					t.Logf("Original: %+v", tt.header)
					t.Logf("Read: %+v", header)
				}
			}
		})
	}
}

// TestPackageHeaderRoundTrip verifies round-trip serialization
func TestPackageHeaderRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		header PackageHeader
	}{
		{
			"Minimal header",
			PackageHeader{
				Magic:              NVPKMagic,
				FormatVersion:      FormatVersion,
				PackageDataVersion: 1,
				MetadataVersion:    1,
				Reserved:           0,
				ArchivePartInfo:    0x00010001,
			},
		},
		{"Full header", packageHeaderRoundTripFull()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write
			var buf bytes.Buffer
			if _, err := tt.header.writeTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var header PackageHeader
			if _, err := header.readFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			// Compare
			if header != tt.header {
				t.Errorf("Round-trip header mismatch")
				t.Logf("Original: %+v", tt.header)
				t.Logf("Round-trip: %+v", header)
			}

			// Validate
			if err := header.validate(); err != nil {
				t.Errorf("Round-trip header validation failed: %v", err)
			}
		})
	}
}

// TestPackageHeaderWriteToErrorPaths verifies WriteTo error handling
//
//nolint:gocognit // table-driven error paths
func TestPackageHeaderWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		header    PackageHeader
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during header write",
			PackageHeader{
				Magic:              NVPKMagic,
				FormatVersion:      FormatVersion,
				PackageDataVersion: 1,
				MetadataVersion:    1,
			},
			testhelpers.NewErrorWriter(),
			true,
			"failed to write header",
		},
		{
			"Failing writer during header write",
			PackageHeader{
				Magic:              NVPKMagic,
				FormatVersion:      FormatVersion,
				PackageDataVersion: 1,
				MetadataVersion:    1,
			},
			testhelpers.NewFailingWriter(50),
			true,
			"failed to write header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.header.writeTo(tt.writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errSubstr != "" {
					errStr := err.Error()
					if !strings.Contains(errStr, tt.errSubstr) {
						t.Errorf("WriteTo() error = %q, want error containing %q", errStr, tt.errSubstr)
					}
				}
			}
		})
	}
}
