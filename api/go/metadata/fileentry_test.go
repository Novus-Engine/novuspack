package metadata

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestFileEntryFixedSize verifies the fixed section is exactly 64 bytes
func TestFileEntryFixedSize(t *testing.T) {
	// Create a minimal FileEntry with only fixed fields
	type FileEntryFixed struct {
		FileID             uint64
		OriginalSize       uint64
		StoredSize         uint64
		RawChecksum        uint32
		StoredChecksum     uint32
		FileVersion        uint32
		MetadataVersion    uint32
		PathCount          uint16
		Type               uint16
		CompressionType    uint8
		CompressionLevel   uint8
		EncryptionType     uint8
		HashCount          uint8
		HashDataOffset     uint32
		HashDataLen        uint16
		OptionalDataLen    uint16
		OptionalDataOffset uint32
		Reserved           uint32
	}

	var fixed FileEntryFixed
	size := binary.Size(fixed)

	if size != FileEntryFixedSize {
		t.Errorf("FileEntry fixed section size = %d bytes, want %d bytes", size, FileEntryFixedSize)
	}
}

// TestFileEntryValidation verifies validation logic
func TestFileEntryValidation(t *testing.T) {
	tests := []struct {
		name    string
		entry   FileEntry
		wantErr bool
	}{
		{
			"Valid entry",
			FileEntry{
				FileID:       1,
				PathCount:    1,
				Paths:        []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				HashCount:    0,
				Hashes:       []HashEntry{},
				OptionalData: []OptionalDataEntry{},
			},
			false,
		},
		{
			"Zero FileID",
			FileEntry{
				FileID:    0,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
			},
			true,
		},
		{
			"Non-zero reserved",
			FileEntry{
				FileID:   1,
				Reserved: 1,
			},
			true,
		},
		{
			"PathCount mismatch",
			FileEntry{
				FileID:    1,
				PathCount: 2,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
			},
			true,
		},
		{
			"HashCount mismatch",
			FileEntry{
				FileID:    1,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				HashCount: 1,
				Hashes:    []HashEntry{},
			},
			true,
		},
		{
			"Invalid hash entry",
			FileEntry{
				FileID:    1,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				HashCount: 1,
				Hashes: []HashEntry{
					{HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification, HashLength: 32, HashData: make([]byte, 16)},
				},
			},
			true,
		},
		{
			"Invalid optional data entry",
			FileEntry{
				FileID:    1,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				OptionalData: []OptionalDataEntry{
					{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 5)},
				},
			},
			true,
		},
		{
			"Valid entry with hashes",
			FileEntry{
				FileID:    1,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				HashCount: 1,
				Hashes: []HashEntry{
					{HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification, HashLength: 32, HashData: make([]byte, 32)},
				},
			},
			false,
		},
		{
			"Valid entry with optional data",
			FileEntry{
				FileID:    1,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				OptionalData: []OptionalDataEntry{
					{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
				},
			},
			false,
		},
		{
			"Valid entry with hashes and optional data",
			FileEntry{
				FileID:    1,
				PathCount: 1,
				Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
				HashCount: 1,
				Hashes: []HashEntry{
					{HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification, HashLength: 32, HashData: make([]byte, 32)},
				},
				OptionalData: []OptionalDataEntry{
					{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFileEntrySizeCalculation verifies size calculations
func TestFileEntrySizeCalculation(t *testing.T) {
	entry := FileEntry{
		FileID:    1,
		PathCount: 1,
		Paths: []generics.PathEntry{
			{PathLength: 8, Path: "test.txt"},
		},
		HashCount: 1,
		Hashes: []HashEntry{
			{HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification, HashLength: 32, HashData: make([]byte, 32)},
		},
		OptionalData: []OptionalDataEntry{},
	}

	// Fixed size should be 64 bytes
	if entry.FixedSize() != 64 {
		t.Errorf("FixedSize() = %d, want 64", entry.FixedSize())
	}

	// Variable size: path(2+8=10) + hash(4+32=36) = 46
	wantVariableSize := 46
	if entry.VariableSize() != wantVariableSize {
		t.Errorf("VariableSize() = %d, want %d", entry.VariableSize(), wantVariableSize)
	}

	// Total size: 64 + 46 = 110
	wantTotalSize := 110
	if entry.TotalSize() != wantTotalSize {
		t.Errorf("TotalSize() = %d, want %d", entry.TotalSize(), wantTotalSize)
	}
}

// TestFileEntryVariableSizeEmptySlices verifies VariableSize with empty slices
func TestFileEntryVariableSizeEmptySlices(t *testing.T) {
	entry := FileEntry{
		FileID:       1,
		PathCount:    0,
		Paths:        []generics.PathEntry{},
		HashCount:    0,
		Hashes:       []HashEntry{},
		OptionalData: []OptionalDataEntry{},
	}

	// Variable size should be 0 when all slices are empty
	if entry.VariableSize() != 0 {
		t.Errorf("VariableSize() = %d, want 0", entry.VariableSize())
	}

	// Total size should be just fixed size
	if entry.TotalSize() != 64 {
		t.Errorf("TotalSize() = %d, want 64", entry.TotalSize())
	}
}

// TestFileEntryVariableSizeWithOptionalData verifies VariableSize includes optional data
func TestFileEntryVariableSizeWithOptionalData(t *testing.T) {
	entry := FileEntry{
		FileID:    1,
		PathCount: 1,
		Paths:     []generics.PathEntry{{PathLength: 8, Path: "test.txt"}},
		HashCount: 0,
		Hashes:    []HashEntry{},
		OptionalData: []OptionalDataEntry{
			{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
		},
	}

	// Variable size should include path + optional data
	// Path: 2 + 8 = 10
	// Optional data: 3 + 10 = 13
	// Total: 23
	wantVariableSize := 23
	if entry.VariableSize() != wantVariableSize {
		t.Errorf("VariableSize() = %d, want %d", entry.VariableSize(), wantVariableSize)
	}
}

// TestFileEntryValidationInvalidPath verifies validation with invalid path entry
func TestFileEntryValidationInvalidPath(t *testing.T) {
	entry := FileEntry{
		FileID:    1,
		PathCount: 1,
		Paths: []generics.PathEntry{
			{PathLength: 8, Path: ""}, // Empty path should fail validation
		},
		HashCount:    0,
		Hashes:       []HashEntry{},
		OptionalData: []OptionalDataEntry{},
	}

	err := entry.Validate()
	if err == nil {
		t.Error("Validate() expected error for invalid path entry, got nil")
	}
}

// TestNewFileEntry verifies NewFileEntry initializes correctly
func TestNewFileEntry(t *testing.T) {
	entry := NewFileEntry()

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so entry is not nil after check
	if entry == nil {
		t.Fatal("NewFileEntry() returned nil")
	}

	// Verify all fields are zero or empty
	//nolint:staticcheck // SA5011: false positive - entry is guaranteed non-nil from NewFileEntry()
	if entry.FileID != 0 {
		t.Errorf("FileID = %d, want 0", entry.FileID)
	}
	if entry.OriginalSize != 0 {
		t.Errorf("OriginalSize = %d, want 0", entry.OriginalSize)
	}
	if entry.StoredSize != 0 {
		t.Errorf("StoredSize = %d, want 0", entry.StoredSize)
	}
	if entry.PathCount != 0 {
		t.Errorf("PathCount = %d, want 0", entry.PathCount)
	}
	if entry.HashCount != 0 {
		t.Errorf("HashCount = %d, want 0", entry.HashCount)
	}
	if len(entry.Paths) != 0 {
		t.Errorf("Paths length = %d, want 0", len(entry.Paths))
	}
	if len(entry.Hashes) != 0 {
		t.Errorf("Hashes length = %d, want 0", len(entry.Hashes))
	}
	if len(entry.OptionalData) != 0 {
		t.Errorf("OptionalData length = %d, want 0", len(entry.OptionalData))
	}

	// Verify runtime fields are initialized correctly
	if entry.ProcessingState != ProcessingStateIdle {
		t.Errorf("ProcessingState = %v, want ProcessingStateIdle", entry.ProcessingState)
	}
	if entry.Data != nil {
		t.Errorf("Data = %v, want nil", entry.Data)
	}
	if entry.SourceFile != nil {
		t.Errorf("SourceFile = %v, want nil", entry.SourceFile)
	}
	if entry.SourceOffset != 0 {
		t.Errorf("SourceOffset = %d, want 0", entry.SourceOffset)
	}
	if entry.SourceSize != 0 {
		t.Errorf("SourceSize = %d, want 0", entry.SourceSize)
	}
	if entry.TempFilePath != "" {
		t.Errorf("TempFilePath = %q, want empty string", entry.TempFilePath)
	}
	if entry.IsDataLoaded {
		t.Errorf("IsDataLoaded = %v, want false", entry.IsDataLoaded)
	}
	if entry.IsTempFile {
		t.Errorf("IsTempFile = %v, want false", entry.IsTempFile)
	}
	if len(entry.PathMetadataEntries) > 0 {
		t.Errorf("PathMetadataEntries = %v, want nil or empty", entry.PathMetadataEntries)
	}
}

// TestFileEntryReadFromFixedOnly verifies ReadFrom with fixed section only
func TestFileEntryReadFromFixedOnly(t *testing.T) {
	entry := FileEntry{
		FileID:           1,
		OriginalSize:     1000,
		StoredSize:       800,
		RawChecksum:      0x12345678,
		StoredChecksum:   0x87654321,
		FileVersion:      1,
		MetadataVersion:  1,
		PathCount:        0,
		Type:             0x0001,
		CompressionType:  fileformat.CompressionZstd,
		CompressionLevel: 6,
		EncryptionType:   fileformat.EncryptionNone,
		HashCount:        0,
		Reserved:         0,
	}
	// Set data so WriteTo can succeed
	entry.SetData([]byte("test data"))

	// Serialize using WriteTo
	buf := new(bytes.Buffer)
	_, err := entry.WriteTo(buf)
	if err != nil {
		t.Fatalf("Failed to serialize entry: %v", err)
	}

	// Deserialize using ReadFrom
	var readEntry FileEntry
	n, err := readEntry.ReadFrom(buf)

	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	if n != int64(FileEntryFixedSize) {
		t.Errorf("ReadFrom() read %d bytes, want %d", n, FileEntryFixedSize)
	}

	// Verify fixed fields match
	if readEntry.FileID != entry.FileID {
		t.Errorf("FileID = %d, want %d", readEntry.FileID, entry.FileID)
	}
	if readEntry.OriginalSize != entry.OriginalSize {
		t.Errorf("OriginalSize = %d, want %d", readEntry.OriginalSize, entry.OriginalSize)
	}
}

// TestFileEntryReadFromWithVariableData verifies ReadFrom with variable-length data
func TestFileEntryReadFromWithVariableData(t *testing.T) {
	entry := FileEntry{
		FileID:          1,
		OriginalSize:    1000,
		StoredSize:      800,
		FileVersion:     1,
		MetadataVersion: 1,
		PathCount:       1,
		HashCount:       1,
		Reserved:        0,
		Paths: []generics.PathEntry{
			{
				PathLength: 8,
				Path:       "test.txt",
			},
		},
		Hashes: []HashEntry{
			{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
		},
	}

	// Update PathLength to match actual path length
	for i := range entry.Paths {
		entry.Paths[i].PathLength = uint16(len(entry.Paths[i].Path))
	}

	// Update HashLength to match actual hash data length
	for i := range entry.Hashes {
		entry.Hashes[i].HashLength = uint16(len(entry.Hashes[i].HashData))
	}
	// Set data so WriteTo can succeed
	entry.SetData([]byte("test data"))

	// Serialize using WriteTo
	buf := new(bytes.Buffer)
	_, writeErr := entry.WriteTo(buf)
	if writeErr != nil {
		t.Fatalf("Failed to serialize entry: %v", writeErr)
	}

	// Deserialize using ReadFrom
	var readEntry FileEntry
	_, readErr := readEntry.ReadFrom(buf)

	if readErr != nil {
		t.Fatalf("ReadFrom() error = %v", readErr)
	}

	// Verify counts
	if readEntry.PathCount != entry.PathCount {
		t.Errorf("PathCount = %d, want %d", readEntry.PathCount, entry.PathCount)
	}
	if readEntry.HashCount != entry.HashCount {
		t.Errorf("HashCount = %d, want %d", readEntry.HashCount, entry.HashCount)
	}
}

// TestFileEntryReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestFileEntryReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Partial fixed", make([]byte, 32)},
		{"Almost complete fixed", make([]byte, 63)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry FileEntry
			r := bytes.NewReader(tt.data)
			_, err := entry.ReadFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestFileEntryWriteToFixedOnly verifies WriteTo with fixed section and data
func TestFileEntryWriteToFixedOnly(t *testing.T) {
	entry := FileEntry{
		FileID:          1,
		OriginalSize:    1000,
		StoredSize:      800,
		FileVersion:     1,
		MetadataVersion: 1,
		PathCount:       0,
		HashCount:       0,
		Reserved:        0,
	}
	// Set data so WriteTo can succeed
	entry.SetData([]byte("test data"))

	var buf bytes.Buffer
	n, err := entry.WriteTo(&buf)

	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	expectedSize := entry.FixedSize() + len(entry.Data)
	if n != int64(expectedSize) {
		t.Errorf("WriteTo() wrote %d bytes, want %d", n, expectedSize)
	}

	if buf.Len() != expectedSize {
		t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), expectedSize)
	}

	// Verify we can read it back
	var readEntry FileEntry
	_, readErr := readEntry.ReadFrom(&buf)
	if readErr != nil {
		t.Errorf("Failed to read back written data: %v", readErr)
	}

	if readEntry.FileID != entry.FileID {
		t.Errorf("FileID mismatch: %d != %d", readEntry.FileID, entry.FileID)
	}
}

// TestFileEntryWriteToWithVariableData verifies WriteTo with variable-length data
func TestFileEntryWriteToWithVariableData(t *testing.T) {
	entry := FileEntry{
		FileID:          1,
		OriginalSize:    1000,
		StoredSize:      800,
		FileVersion:     1,
		MetadataVersion: 1,
		PathCount:       1,
		HashCount:       1,
		Reserved:        0,
		Paths: []generics.PathEntry{
			{
				PathLength: 8,
				Path:       "test.txt",
			},
		},
		Hashes: []HashEntry{
			{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
		},
	}

	// Update counts
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	// Set data so WriteTo can succeed
	entry.SetData([]byte("test data"))

	var buf bytes.Buffer
	_, err := entry.WriteTo(&buf)

	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	// Verify we can read it back
	var readEntry FileEntry
	_, readErr := readEntry.ReadFrom(&buf)
	if readErr != nil {
		t.Fatalf("Failed to read back written data: %v", readErr)
	}

	// Verify counts
	if readEntry.PathCount != entry.PathCount {
		t.Errorf("PathCount mismatch: %d != %d", readEntry.PathCount, entry.PathCount)
	}
	if readEntry.HashCount != entry.HashCount {
		t.Errorf("HashCount mismatch: %d != %d", readEntry.HashCount, entry.HashCount)
	}
}

// TestFileEntryRoundTrip verifies round-trip serialization
func TestFileEntryRoundTrip(t *testing.T) {
	entry := FileEntry{
		FileID:           1,
		OriginalSize:     10000,
		StoredSize:       8000,
		RawChecksum:      0x12345678,
		StoredChecksum:   0x87654321,
		FileVersion:      5,
		MetadataVersion:  3,
		PathCount:        2,
		Type:             0x0001,
		CompressionType:  fileformat.CompressionZstd,
		CompressionLevel: 6,
		EncryptionType:   fileformat.EncryptionNone,
		HashCount:        2,
		Reserved:         0,
		Paths: []generics.PathEntry{
			{PathLength: 8, Path: "test.txt"},
			{PathLength: 9, Path: "test2.txt"},
		},
		Hashes: []HashEntry{
			{HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification, HashLength: 32, HashData: make([]byte, 32)},
			{HashType: fileformat.HashTypeSHA512, HashPurpose: fileformat.HashPurposeDeduplication, HashLength: 64, HashData: make([]byte, 64)},
		},
		OptionalData: []OptionalDataEntry{
			{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
		},
	}

	// Update counts
	entry.PathCount = uint16(len(entry.Paths))
	entry.HashCount = uint8(len(entry.Hashes))
	// Set data so WriteTo can succeed
	entry.SetData([]byte("test data"))

	// Write
	var buf bytes.Buffer
	if _, err := entry.WriteTo(&buf); err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	// Read
	var readEntry FileEntry
	if _, err := readEntry.ReadFrom(&buf); err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	// Compare fixed fields
	if readEntry.FileID != entry.FileID {
		t.Errorf("FileID mismatch: %d != %d", readEntry.FileID, entry.FileID)
	}
	if readEntry.OriginalSize != entry.OriginalSize {
		t.Errorf("OriginalSize mismatch: %d != %d", readEntry.OriginalSize, entry.OriginalSize)
	}

	// Compare counts
	if readEntry.PathCount != entry.PathCount {
		t.Errorf("PathCount mismatch: %d != %d", readEntry.PathCount, entry.PathCount)
	}
	if readEntry.HashCount != entry.HashCount {
		t.Errorf("HashCount mismatch: %d != %d", readEntry.HashCount, entry.HashCount)
	}

	// Validate
	if err := readEntry.Validate(); err != nil {
		t.Errorf("Round-trip entry validation failed: %v", err)
	}
}

// TestFileEntryWriteToErrorPaths verifies WriteTo error handling
func TestFileEntryWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		entry     FileEntry
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during fixed section write",
			FileEntry{
				FileID: 1,
				Paths: []generics.PathEntry{
					{PathLength: 8, Path: "test.txt"},
				},
			},
			testhelpers.NewErrorWriter(),
			true,
			"failed to write",
		},
		{
			"Failing writer during path write",
			FileEntry{
				FileID: 1,
				Paths: []generics.PathEntry{
					{PathLength: 8, Path: "test.txt"},
				},
			},
			testhelpers.NewFailingWriter(63),
			true,
			"failed to write",
		},
		{
			"Failing writer during hash write",
			FileEntry{
				FileID: 1,
				Paths: []generics.PathEntry{
					{PathLength: 8, Path: "test.txt"},
				},
				Hashes: []HashEntry{
					{HashType: fileformat.HashTypeSHA256, HashLength: 32, HashData: make([]byte, 32)},
				},
				Data: []byte("test data"), // Set data so WriteTo can proceed
			},
			testhelpers.NewFailingWriter(50), // Fail during hash write in metadata section
			true,
			"failed to write",
		},
		{
			"Path length mismatch",
			FileEntry{
				FileID: 1,
				Paths: []generics.PathEntry{
					{PathLength: 20, Path: "test.txt"}, // 8 bytes, but PathLength says 20
				},
			},
			&bytes.Buffer{},
			true,
			"failed to write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update counts and lengths
			tt.entry.PathCount = uint16(len(tt.entry.Paths))
			tt.entry.HashCount = uint8(len(tt.entry.Hashes))
			for i := range tt.entry.Paths {
				if !strings.Contains(tt.name, "mismatch") {
					tt.entry.Paths[i].PathLength = uint16(len(tt.entry.Paths[i].Path))
				}
			}
			for i := range tt.entry.Hashes {
				tt.entry.Hashes[i].HashLength = uint16(len(tt.entry.Hashes[i].HashData))
			}
			_, err := tt.entry.WriteTo(tt.writer)

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

// TestFileEntryReadFromErrorPaths verifies ReadFrom error handling for edge cases
func TestFileEntryReadFromErrorPaths(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *bytes.Reader
		wantErr bool
	}{
		{
			"Reader with no data",
			func() *bytes.Reader {
				return bytes.NewReader([]byte{})
			},
			true,
		},
		{
			"Reader with partial fixed section",
			func() *bytes.Reader {
				return bytes.NewReader(make([]byte, 32))
			},
			true,
		},
		{
			"Reader with fixed section but incomplete path",
			func() *bytes.Reader {
				// Create a FileEntry with one path, but truncate the path data
				entry := FileEntry{
					FileID:    1,
					PathCount: 1,
					Paths: []generics.PathEntry{
						{PathLength: 8, Path: "test.txt"},
					},
				}
				entry.SetData([]byte("test data"))
				var buf bytes.Buffer
				if _, err := entry.WriteTo(&buf); err != nil {
					panic(err)
				}
				// Truncate to remove part of the path
				data := buf.Bytes()
				return bytes.NewReader(data[:64+2+4]) // Fixed section + PathLength + partial path
			},
			true,
		},
		{
			"Reader with paths but incomplete hash when HashCount > 0",
			func() *bytes.Reader {
				entry := FileEntry{
					FileID:    1,
					PathCount: 1,
					HashCount: 1,
					Paths: []generics.PathEntry{
						{PathLength: 8, Path: "test.txt"},
					},
					Hashes: []HashEntry{
						{HashType: fileformat.HashTypeSHA256, HashLength: 32, HashData: make([]byte, 32)},
					},
				}
				entry.SetData([]byte("test data"))
				var buf bytes.Buffer
				if _, err := entry.WriteTo(&buf); err != nil {
					panic(err)
				}
				// Truncate to remove part of the hash data
				data := buf.Bytes()
				return bytes.NewReader(data[:76]) // Fixed (64) + path (10) + partial hash header (2 of 4 bytes)
			},
			true,
		},
		{
			"Reader with paths and hashes but incomplete optional data when OptionalDataLen > 0",
			func() *bytes.Reader {
				entry := FileEntry{
					FileID:    1,
					PathCount: 1,
					HashCount: 1,
					Paths: []generics.PathEntry{
						{PathLength: 8, Path: "test.txt"},
					},
					Hashes: []HashEntry{
						{HashType: fileformat.HashTypeSHA256, HashLength: 32, HashData: make([]byte, 32)},
					},
					OptionalData: []OptionalDataEntry{
						{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
					},
				}
				entry.SetData([]byte("test data"))
				var buf bytes.Buffer
				if _, err := entry.WriteTo(&buf); err != nil {
					panic(err)
				}
				// Truncate to remove part of the optional data (before optional data is fully read)
				data := buf.Bytes()
				// Truncate after fixed + path + hash, but before completing optional data read
				// Fixed (64) + path (10) + hash (36) + optional header start (2 bytes of 3-byte header) = 112
				return bytes.NewReader(data[:112])
			},
			true,
		},
		{
			"Reader with HashDataOffset requiring skip",
			func() *bytes.Reader {
				entry := FileEntry{
					FileID:         1,
					PathCount:      1,
					HashCount:      1,
					HashDataOffset: 100, // Larger than path size, requires skip
					Paths: []generics.PathEntry{
						{PathLength: 8, Path: "test.txt"},
					},
					Hashes: []HashEntry{
						{HashType: fileformat.HashTypeSHA256, HashLength: 32, HashData: make([]byte, 32)},
					},
				}
				var buf bytes.Buffer
				// Manually create data with skip
				_ = binary.Write(&buf, binary.LittleEndian, entry.FileID)
				_ = binary.Write(&buf, binary.LittleEndian, entry.OriginalSize)
				_ = binary.Write(&buf, binary.LittleEndian, entry.StoredSize)
				_ = binary.Write(&buf, binary.LittleEndian, entry.RawChecksum)
				_ = binary.Write(&buf, binary.LittleEndian, entry.StoredChecksum)
				_ = binary.Write(&buf, binary.LittleEndian, entry.FileVersion)
				_ = binary.Write(&buf, binary.LittleEndian, entry.MetadataVersion)
				_ = binary.Write(&buf, binary.LittleEndian, uint16(1)) // PathCount
				_ = binary.Write(&buf, binary.LittleEndian, entry.Type)
				_ = binary.Write(&buf, binary.LittleEndian, entry.CompressionType)
				_ = binary.Write(&buf, binary.LittleEndian, entry.CompressionLevel)
				_ = binary.Write(&buf, binary.LittleEndian, entry.EncryptionType)
				_ = binary.Write(&buf, binary.LittleEndian, uint8(1))    // HashCount
				_ = binary.Write(&buf, binary.LittleEndian, uint32(100)) // HashDataOffset
				_ = binary.Write(&buf, binary.LittleEndian, uint16(36))  // HashDataLen
				_ = binary.Write(&buf, binary.LittleEndian, uint16(0))   // OptionalDataLen
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // OptionalDataOffset
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // Reserved
				// Write path
				if _, err := entry.Paths[0].WriteTo(&buf); err != nil {
					panic(err)
				}
				// Write padding to reach HashDataOffset (path is 10 bytes, need 90 more to reach offset 100)
				padding := make([]byte, 100-10)
				buf.Write(padding)
				// Write hash
				if _, err := entry.Hashes[0].WriteTo(&buf); err != nil {
					panic(err)
				}
				return bytes.NewReader(buf.Bytes())
			},
			false, // Should succeed
		},
		{
			"Reader with HashDataOffset skip failure",
			func() *bytes.Reader {
				// Manually create data with HashDataOffset requiring skip but not enough bytes
				var buf bytes.Buffer
				// Write fixed section with HashDataOffset=200
				_ = binary.Write(&buf, binary.LittleEndian, uint64(1))   // FileID
				_ = binary.Write(&buf, binary.LittleEndian, uint64(0))   // OriginalSize
				_ = binary.Write(&buf, binary.LittleEndian, uint64(0))   // StoredSize
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // RawChecksum
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // StoredChecksum
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // FileVersion
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // MetadataVersion
				_ = binary.Write(&buf, binary.LittleEndian, uint16(1))   // PathCount
				_ = binary.Write(&buf, binary.LittleEndian, uint16(0))   // Type
				_ = binary.Write(&buf, binary.LittleEndian, uint8(0))    // CompressionType
				_ = binary.Write(&buf, binary.LittleEndian, uint8(0))    // CompressionLevel
				_ = binary.Write(&buf, binary.LittleEndian, uint8(0))    // EncryptionType
				_ = binary.Write(&buf, binary.LittleEndian, uint8(1))    // HashCount
				_ = binary.Write(&buf, binary.LittleEndian, uint32(100)) // HashDataOffset (requires skip from path end at 10)
				_ = binary.Write(&buf, binary.LittleEndian, uint16(36))  // HashDataLen
				_ = binary.Write(&buf, binary.LittleEndian, uint16(0))   // OptionalDataLen
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // OptionalDataOffset
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // Reserved
				// Write path (10 bytes: 2 PathLength + 8 Path)
				pathEntry := generics.PathEntry{PathLength: 8, Path: "test.txt"}
				if _, err := pathEntry.WriteTo(&buf); err != nil {
					panic(err)
				}
				// Only write 50 bytes of padding (need 90 more to reach hash at offset 100, but only write 50)
				padding := make([]byte, 50)
				buf.Write(padding)
				// Truncate here - not enough data to reach hash (would need 90 bytes but only have 50)
				return bytes.NewReader(buf.Bytes())
			},
			true,
		},
		{
			"Reader with OptionalDataOffset requiring skip",
			func() *bytes.Reader {
				entry := FileEntry{
					FileID:             1,
					PathCount:          1,
					HashCount:          1,
					OptionalDataOffset: 200, // Requires skip
					OptionalDataLen:    10,
					Paths: []generics.PathEntry{
						{PathLength: 8, Path: "test.txt"},
					},
					Hashes: []HashEntry{
						{HashType: fileformat.HashTypeSHA256, HashLength: 32, HashData: make([]byte, 32)},
					},
					OptionalData: []OptionalDataEntry{
						{DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10)},
					},
				}
				entry.SetData([]byte("test data"))
				var buf bytes.Buffer
				if _, err := entry.WriteTo(&buf); err != nil {
					panic(err)
				}
				return bytes.NewReader(buf.Bytes())
			},
			false, // Should succeed
		},
		{
			"Reader with OptionalDataOffset skip failure",
			func() *bytes.Reader {
				// Manually create data with OptionalDataOffset requiring skip but not enough bytes
				var buf bytes.Buffer
				// Write fixed section with OptionalDataOffset=300
				_ = binary.Write(&buf, binary.LittleEndian, uint64(1))   // FileID
				_ = binary.Write(&buf, binary.LittleEndian, uint64(0))   // OriginalSize
				_ = binary.Write(&buf, binary.LittleEndian, uint64(0))   // StoredSize
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // RawChecksum
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // StoredChecksum
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // FileVersion
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // MetadataVersion
				_ = binary.Write(&buf, binary.LittleEndian, uint16(1))   // PathCount
				_ = binary.Write(&buf, binary.LittleEndian, uint16(0))   // Type
				_ = binary.Write(&buf, binary.LittleEndian, uint8(0))    // CompressionType
				_ = binary.Write(&buf, binary.LittleEndian, uint8(0))    // CompressionLevel
				_ = binary.Write(&buf, binary.LittleEndian, uint8(0))    // EncryptionType
				_ = binary.Write(&buf, binary.LittleEndian, uint8(1))    // HashCount
				_ = binary.Write(&buf, binary.LittleEndian, uint32(10))  // HashDataOffset (after 10-byte path)
				_ = binary.Write(&buf, binary.LittleEndian, uint16(36))  // HashDataLen
				_ = binary.Write(&buf, binary.LittleEndian, uint16(10))  // OptionalDataLen
				_ = binary.Write(&buf, binary.LittleEndian, uint32(300)) // OptionalDataOffset (requires skip)
				_ = binary.Write(&buf, binary.LittleEndian, uint32(0))   // Reserved
				// Write path (10 bytes: 2 PathLength + 8 Path)
				pathEntry := generics.PathEntry{PathLength: 8, Path: "test.txt"}
				if _, err := pathEntry.WriteTo(&buf); err != nil {
					panic(err)
				}
				// Write hash (36 bytes)
				hashEntry := HashEntry{HashType: fileformat.HashTypeSHA256, HashLength: 32, HashData: make([]byte, 32)}
				if _, err := hashEntry.WriteTo(&buf); err != nil {
					panic(err)
				}
				// Only write 50 bytes of padding (need more to reach optional data at offset 300, but only write 50)
				padding := make([]byte, 50)
				buf.Write(padding)
				// Truncate here - not enough data to reach optional data
				return bytes.NewReader(buf.Bytes())
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry FileEntry
			r := tt.setup()
			_, err := entry.ReadFrom(r)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
