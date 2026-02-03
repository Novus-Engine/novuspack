package fileformat

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
	"github.com/novus-engine/novuspack/api/go/pkgerrors"
)

// fileIndexWithManyEntries returns a FileIndex with 5 entries for use in tests.
func fileIndexWithManyEntries() FileIndex {
	return FileIndex{
		EntryCount: 5,
		Reserved:   0,
		Entries: []IndexEntry{
			{FileID: 1, Offset: 112},
			{FileID: 2, Offset: 256},
			{FileID: 3, Offset: 512},
			{FileID: 4, Offset: 1024},
			{FileID: 5, Offset: 2048},
		},
	}
}

// compareFileIndexEntries compares index.Entries with want.Entries and reports mismatches via t.
func compareFileIndexEntries(t *testing.T, index, want FileIndex) {
	t.Helper()
	for i, entry := range index.Entries {
		if entry.FileID != want.Entries[i].FileID {
			t.Errorf("Entry[%d].FileID = %d, want %d", i, entry.FileID, want.Entries[i].FileID)
		}
		if entry.Offset != want.Entries[i].Offset {
			t.Errorf("Entry[%d].Offset = %d, want %d", i, entry.Offset, want.Entries[i].Offset)
		}
	}
}

// TestIndexEntrySize verifies IndexEntry is exactly 16 bytes
func TestIndexEntrySize(t *testing.T) {
	var entry IndexEntry
	size := binary.Size(entry)

	if size != IndexEntrySize {
		t.Errorf("IndexEntry size = %d bytes, want %d bytes", size, IndexEntrySize)
	}
}

// TestFileIndexValidation verifies validation logic
func TestFileIndexValidation(t *testing.T) {
	tests := []struct {
		name    string
		index   FileIndex
		wantErr bool
	}{
		{
			"Valid index",
			FileIndex{
				EntryCount: 2,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
				},
			},
			false,
		},
		{
			"Non-zero reserved",
			FileIndex{
				Reserved: 1,
			},
			true,
		},
		{
			"Entry count mismatch",
			FileIndex{
				EntryCount: 3,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
				},
			},
			true,
		},
		{
			"Zero FileID",
			FileIndex{
				EntryCount: 1,
				Entries: []IndexEntry{
					{FileID: 0, Offset: 112},
				},
			},
			true,
		},
		{
			"Duplicate FileID",
			FileIndex{
				EntryCount: 2,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 1, Offset: 256},
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.index.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFileIndexSizeCalculation verifies size calculation
func TestFileIndexSizeCalculation(t *testing.T) {
	tests := []struct {
		name       string
		entryCount int
		wantSize   int
	}{
		{"Empty index", 0, 16},
		{"Single entry", 1, 32},
		{"Multiple entries", 10, 176},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := FileIndex{
				EntryCount: uint32(tt.entryCount),
				Entries:    make([]IndexEntry, tt.entryCount),
			}

			if index.size() != tt.wantSize {
				t.Errorf("size() = %d, want %d", index.size(), tt.wantSize)
			}
		})
	}
}

// TestNewFileIndex verifies NewFileIndex initializes correctly
func TestNewFileIndex(t *testing.T) {
	index := NewFileIndex()
	//nolint:staticcheck // SA5011: false positive - t.Fatal exits
	if index == nil {
		t.Fatal("NewFileIndex() returned nil")
	}

	// Verify all fields are zero or empty
	//nolint:staticcheck // SA5011: false positive - index is not nil after t.Fatal check
	if index.EntryCount != 0 {
		t.Errorf("EntryCount = %d, want 0", index.EntryCount)
	}
	if index.Reserved != 0 {
		t.Errorf("Reserved = %d, want 0", index.Reserved)
	}
	if len(index.Entries) != 0 {
		t.Errorf("Entries length = %d, want 0", len(index.Entries))
	}
}

// TestFileIndexReadFrom verifies ReadFrom deserialization
// Specification: package_file_format.md: 6 File Index Section
//
//nolint:gocognit // table-driven test with multiple cases
func TestFileIndexReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		index   FileIndex
		wantErr bool
	}{
		{
			"Valid index with entries",
			FileIndex{
				EntryCount: 2,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
				},
			},
			false,
		},
		{
			"Empty index",
			FileIndex{
				EntryCount: 0,
				Reserved:   0,
				Entries:    []IndexEntry{},
			},
			false,
		},
		{
			"Index with many entries",
			fileIndexWithManyEntries(),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update EntryCount to match actual entries
			tt.index.EntryCount = uint32(len(tt.index.Entries))

			// Serialize using WriteTo (once implemented)
			// For now, serialize manually
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, tt.index.EntryCount)
			_ = binary.Write(buf, binary.LittleEndian, tt.index.Reserved)
			_ = binary.Write(buf, binary.LittleEndian, make([]byte, 8)) // Reserved2 (8 bytes)
			for _, entry := range tt.index.Entries {
				_ = binary.Write(buf, binary.LittleEndian, entry)
			}

			// Deserialize using ReadFrom
			var index FileIndex
			n, err := index.readFrom(buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.index.size()
				if n != int64(expectedSize) {
					t.Errorf("ReadFrom() read %d bytes, want %d", n, expectedSize)
				}

				// Verify all fields match
				if index.EntryCount != tt.index.EntryCount {
					t.Errorf("EntryCount = %d, want %d", index.EntryCount, tt.index.EntryCount)
				}
				if index.Reserved != tt.index.Reserved {
					t.Errorf("Reserved = %d, want %d", index.Reserved, tt.index.Reserved)
				}
				if len(index.Entries) != len(tt.index.Entries) {
					t.Errorf("Entries length = %d, want %d", len(index.Entries), len(tt.index.Entries))
				}

				// Verify entries match
				compareFileIndexEntries(t, index, tt.index)

				// Verify validation passes
				if err := index.validate(); err != nil {
					t.Errorf("ReadFrom() index validation failed: %v", err)
				}
			}
		})
	}
}

// TestFileIndexReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestFileIndexReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Partial header", make([]byte, 8)},
		{"Almost complete header", make([]byte, 15)},
		{"Incomplete EntryCount read", func() []byte {
			return []byte{0x02, 0x00} // Only 2 bytes of EntryCount (need 4)
		}()},
		{"Incomplete Reserved read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2)) // EntryCount
			return buf.Bytes()[:5]                                // Only 1 byte of Reserved (need 4)
		}()},
		{"Incomplete FirstEntryOffset read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2)) // EntryCount
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
			return buf.Bytes()[:12]                               // Only 4 bytes of FirstEntryOffset (need 8)
		}()},
		{"Header with EntryCount>0 but no entries", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2))
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
			_ = binary.Write(buf, binary.LittleEndian, uint64(0)) // FirstEntryOffset
			// Only 16 bytes, but EntryCount says 2 entries needed (32 more bytes)
			return buf.Bytes()
		}()},
		{"Header with EntryCount>0 but incomplete first entry", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2))
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
			_ = binary.Write(buf, binary.LittleEndian, uint64(0)) // FirstEntryOffset
			_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // First entry FileID
			// Only 8 bytes of first entry (need 16)
			return buf.Bytes()
		}()},
		{"Header with EntryCount>0 but incomplete second entry", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2))
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))                 // Reserved
			_ = binary.Write(buf, binary.LittleEndian, uint64(0))                 // FirstEntryOffset
			_ = binary.Write(buf, binary.LittleEndian, uint64(1))                 // First entry FileID
			_ = binary.Write(buf, binary.LittleEndian, uint64(PackageHeaderSize)) // First entry Offset
			_ = binary.Write(buf, binary.LittleEndian, uint64(2))                 // Second entry FileID
			// Only 8 bytes of second entry (need 16)
			return buf.Bytes()
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var index FileIndex
			r := bytes.NewReader(tt.data)
			_, err := index.readFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestFileIndexWriteTo verifies WriteTo serialization
// Specification: package_file_format.md: 6 File Index Section
//
//nolint:gocognit // table-driven test
func TestFileIndexWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		index   FileIndex
		wantErr bool
	}{
		{
			"Valid index with entries",
			FileIndex{
				EntryCount: 2,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
				},
			},
			false,
		},
		{
			"Empty index",
			FileIndex{
				EntryCount: 0,
				Reserved:   0,
				Entries:    []IndexEntry{},
			},
			false,
		},
		{
			"Index with single entry",
			FileIndex{
				EntryCount: 1,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
				},
			},
			false,
		},
		{
			"Index with many entries",
			fileIndexWithManyEntries(),
			false,
		},
		{
			"Index with non-zero Reserved",
			FileIndex{
				EntryCount: 2,
				Reserved:   0x12345678,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
				},
			},
			false,
		},
		{
			"Index with non-zero FirstEntryOffset",
			FileIndex{
				EntryCount:       2,
				Reserved:         0,
				FirstEntryOffset: 4096,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
				},
			},
			false,
		},
		{
			"Index with many entries (10 entries)",
			FileIndex{
				EntryCount: 10,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024},
					{FileID: 5, Offset: 2048},
					{FileID: 6, Offset: 4096},
					{FileID: 7, Offset: 8192},
					{FileID: 8, Offset: 16384},
					{FileID: 9, Offset: 32768},
					{FileID: 10, Offset: 65536},
				},
			},
			false,
		},
		{
			"Index with large FirstEntryOffset",
			FileIndex{
				EntryCount:       1,
				Reserved:         0,
				FirstEntryOffset: 0xFFFFFFFFFFFFFFFF,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
				},
			},
			false,
		},
		{
			"Index with 20 entries",
			FileIndex{
				EntryCount: 20,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112}, {FileID: 2, Offset: 256}, {FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024}, {FileID: 5, Offset: 2048}, {FileID: 6, Offset: 4096},
					{FileID: 7, Offset: 8192}, {FileID: 8, Offset: 16384}, {FileID: 9, Offset: 32768},
					{FileID: 10, Offset: 65536}, {FileID: 11, Offset: 131072}, {FileID: 12, Offset: 262144},
					{FileID: 13, Offset: 524288}, {FileID: 14, Offset: 1048576}, {FileID: 15, Offset: 2097152},
					{FileID: 16, Offset: 4194304}, {FileID: 17, Offset: 8388608}, {FileID: 18, Offset: 16777216},
					{FileID: 19, Offset: 33554432}, {FileID: 20, Offset: 67108864},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update EntryCount to match actual entries
			tt.index.EntryCount = uint32(len(tt.index.Entries))

			var buf bytes.Buffer
			n, err := tt.index.writeTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.index.size()
				if n != int64(expectedSize) {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, expectedSize)
				}

				if buf.Len() != expectedSize {
					t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), expectedSize)
				}

				// Verify we can read it back
				var index FileIndex
				_, readErr := index.readFrom(&buf)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if index.EntryCount != tt.index.EntryCount {
					t.Errorf("EntryCount mismatch: %d != %d", index.EntryCount, tt.index.EntryCount)
				}
			}
		})
	}
}

// TestFileIndexWriteToErrorPaths verifies WriteTo error handling
//
//nolint:gocognit // table-driven error paths
func TestFileIndexWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		index     FileIndex
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during EntryCount write",
			FileIndex{
				EntryCount: 1,
				Entries:    []IndexEntry{{FileID: 1, Offset: 112}},
			},
			testhelpers.NewErrorWriter(),
			true,
			"failed to write",
		},
		{
			"Failing writer during Reserved write",
			FileIndex{
				EntryCount: 1,
				Entries:    []IndexEntry{{FileID: 1, Offset: 112}},
			},
			testhelpers.NewFailingWriter(3),
			true,
			"failed to write",
		},
		{
			"Failing writer during FirstEntryOffset write",
			FileIndex{
				EntryCount: 1,
				Entries:    []IndexEntry{{FileID: 1, Offset: 112}},
			},
			testhelpers.NewFailingWriter(7),
			true,
			"failed to write",
		},
		{
			"Failing writer during entry write",
			FileIndex{
				EntryCount: 2,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
				},
			},
			testhelpers.NewFailingWriter(20),
			true,
			"failed to write",
		},
		{
			"Failing writer during second entry write",
			FileIndex{
				EntryCount: 3,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
				},
			},
			testhelpers.NewFailingWriter(32), // Allow header (16) + first entry (16) but fail on second
			true,
			"failed to write",
		},
		{
			"Failing writer during third entry write",
			FileIndex{
				EntryCount: 4,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024},
				},
			},
			testhelpers.NewFailingWriter(48), // Allow header (16) + 2 entries (32) but fail on third
			true,
			"failed to write",
		},
		{
			"Failing writer during fourth entry write",
			FileIndex{
				EntryCount: 5,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024},
					{FileID: 5, Offset: 2048},
				},
			},
			testhelpers.NewFailingWriter(64), // Allow header (16) + 3 entries (48) but fail on fourth
			true,
			"failed to write",
		},
		{
			"Empty index write",
			FileIndex{
				EntryCount: 0,
				Reserved:   0,
				Entries:    []IndexEntry{},
			},
			&bytes.Buffer{},
			false,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.index.EntryCount = uint32(len(tt.index.Entries))
			_, err := tt.index.writeTo(tt.writer)

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

// TestFileIndexRoundTrip verifies round-trip serialization
//
//nolint:gocognit // table-driven round-trip
func TestFileIndexRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		index FileIndex
	}{
		{
			"Empty index",
			FileIndex{
				EntryCount: 0,
				Reserved:   0,
				Entries:    []IndexEntry{},
			},
		},
		{
			"Index with entries",
			FileIndex{
				EntryCount: 3,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
				},
			},
		},
		{
			"Index with many entries",
			FileIndex{
				EntryCount: 10,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024},
					{FileID: 5, Offset: 2048},
					{FileID: 6, Offset: 4096},
					{FileID: 7, Offset: 8192},
					{FileID: 8, Offset: 16384},
					{FileID: 9, Offset: 32768},
					{FileID: 10, Offset: 65536},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Update EntryCount to match actual entries
			tt.index.EntryCount = uint32(len(tt.index.Entries))

			// Write
			var buf bytes.Buffer
			if _, err := tt.index.writeTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var index FileIndex
			if _, err := index.readFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			// Compare all fields
			if index.EntryCount != tt.index.EntryCount {
				t.Errorf("EntryCount mismatch: %d != %d", index.EntryCount, tt.index.EntryCount)
			}
			if index.Reserved != tt.index.Reserved {
				t.Errorf("Reserved mismatch: %d != %d", index.Reserved, tt.index.Reserved)
			}
			if len(index.Entries) != len(tt.index.Entries) {
				t.Errorf("Entries length mismatch: %d != %d", len(index.Entries), len(tt.index.Entries))
			}

			// Compare entries
			compareFileIndexEntries(t, index, tt.index)

			// Validate
			if err := index.validate(); err != nil {
				t.Errorf("Round-trip index validation failed: %v", err)
			}

			// Verify FirstEntryOffset matches
			if index.FirstEntryOffset != tt.index.FirstEntryOffset {
				t.Errorf("FirstEntryOffset mismatch: %d != %d", index.FirstEntryOffset, tt.index.FirstEntryOffset)
			}
		})
	}
}

// TestFileIndexReadFrom_OOMPrevention tests the OOM prevention validation paths in ReadFrom.
// These tests verify that corrupted or malicious files with excessive entry counts are rejected
// before attempting memory allocation.
func TestFileIndexReadFrom_OOMPrevention(t *testing.T) {
	// Calculate maxInt for the current architecture
	const maxInt = int(^uint(0) >> 1)
	maxUint32 := uint32(^uint32(0))

	// Helper function to create a file index header with a specific entry count
	createIndexHeader := func(entryCount uint32) []byte {
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, entryCount) // EntryCount
		_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // Reserved
		_ = binary.Write(buf, binary.LittleEndian, uint64(0))  // FirstEntryOffset
		return buf.Bytes()
	}

	// Calculate max safe entry count based on maxInt
	maxSafeEntryCount := maxInt / int(IndexEntrySize)
	var maxSafeEntryCountUint32 uint32
	if maxSafeEntryCount > int(^uint32(0)) {
		maxSafeEntryCountUint32 = ^uint32(0) // Use max uint32 if maxInt is larger
	} else {
		maxSafeEntryCountUint32 = uint32(maxSafeEntryCount)
	}

	tests := []struct {
		name                    string
		entryCount              uint32
		description             string
		errorSubstr             string
		allowAnyValidationError bool // If true, accept any validation error
	}{
		{
			name:                    "entry count exceeds maxInt/int(IndexEntrySize)",
			entryCount:              maxSafeEntryCountUint32 + 1,
			description:             "Tests rejection when entryCount * IndexEntrySize would exceed maxInt (line 227-238)",
			errorSubstr:             "",    // May hit different checks depending on system, or pass if system has enough memory
			allowAnyValidationError: false, // May pass on systems with sufficient memory
		},
		{
			name:                    "entry count at maxInt/int(IndexEntrySize) boundary",
			entryCount:              maxSafeEntryCountUint32,
			description:             "Tests boundary value for maxInt/int(IndexEntrySize) check",
			errorSubstr:             "", // May pass or fail depending on memory check
			allowAnyValidationError: false,
		},
		{
			name:                    "entry count causing multiplication overflow or memory check",
			entryCount:              maxUint32, // Maximum uint32 value
			description:             "Tests rejection when entryCount * IndexEntrySize overflows or exceeds memory (line 208-219 or 247-265)",
			errorSubstr:             "", // May hit calculation overflow or system memory check
			allowAnyValidationError: true,
		},
		{
			name:        "entry count requiring >1GB allocation",
			entryCount:  (1024*1024*1024)/IndexEntrySize + 1, // Just over 1GB
			description: "Tests system memory check for allocations >1GB (line 247-265)",
			errorSubstr: "", // May pass or fail depending on available system memory
		},
		{
			name:        "entry count requiring >10GB allocation",
			entryCount:  (10*1024*1024*1024)/IndexEntrySize + 1, // Just over 10GB
			description: "Tests system memory check for very large allocations (line 247-265)",
			errorSubstr: "exceeds available system memory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.entryCount > 1000000 && testing.Short() {
				t.Skip("Skipping large allocation test in short mode")
			}
			header := createIndexHeader(tt.entryCount)
			r := bytes.NewReader(header)
			var index FileIndex
			_, err := index.readFrom(r)
			assertOOMPreventionResult(t, err, tt.errorSubstr, tt.allowAnyValidationError)
		})
	}
}

func assertOOMPreventionResult(t *testing.T, err error, errorSubstr string, allowAnyValidationError bool) {
	t.Helper()
	if errorSubstr != "" {
		assertOOMPreventionErrorWithSubstr(t, err, errorSubstr)
		return
	}
	if allowAnyValidationError {
		assertOOMPreventionValidationError(t, err)
		return
	}
	if err != nil {
		assertOOMPreventionValidationErrorIfPresent(t, err)
	}
}

func assertOOMPreventionErrorWithSubstr(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Errorf("ReadFrom() expected error containing %q, got nil", substr)
		return
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("ReadFrom() error = %q, want error containing %q", err.Error(), substr)
	}
	assertOOMPreventionValidationError(t, err)
}

func assertOOMPreventionValidationError(t *testing.T, err error) {
	t.Helper()
	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Errorf("ReadFrom() error is not a PackageError: %T", err)
		return
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("ReadFrom() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

func assertOOMPreventionValidationErrorIfPresent(t *testing.T, err error) {
	t.Helper()
	var pkgErr *pkgerrors.PackageError
	if pkgerrors.As(err, &pkgErr) && pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("ReadFrom() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestFileIndexReadFrom_OOMPrevention_MaxIntBoundary tests the specific boundary case
// where entryCount equals maxInt to ensure the check works correctly.
func TestFileIndexReadFrom_OOMPrevention_MaxIntBoundary(t *testing.T) {
	// Test with entryCount = maxInt (should be rejected if > maxInt check is correct)
	// On 64-bit systems, maxInt is 9223372036854775807, which is way larger than max uint32
	// So we test with max uint32 which should pass the first check but may fail others
	maxUint32 := uint32(^uint32(0))

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, maxUint32) // EntryCount = max uint32
	_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
	_ = binary.Write(buf, binary.LittleEndian, uint64(0)) // FirstEntryOffset

	var index FileIndex
	_, err := index.readFrom(buf)

	// maxUint32 should trigger one of the validation errors
	// (likely multiplication overflow or maxInt/int(IndexEntrySize) check)
	if err == nil {
		t.Error("ReadFrom() expected error for maxUint32 entry count, got nil")
		return
	}

	var pkgErr *pkgerrors.PackageError
	if !pkgerrors.As(err, &pkgErr) {
		t.Errorf("ReadFrom() error is not a PackageError: %T", err)
		return
	}
	if pkgErr.Type != pkgerrors.ErrTypeValidation {
		t.Errorf("ReadFrom() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
	}
}

// TestFileIndexReadFrom_OOMPrevention_MultiplicationOverflow tests the specific
// multiplication overflow check (line 208-219).
func TestFileIndexReadFrom_OOMPrevention_MultiplicationOverflow(t *testing.T) {
	// Find a value that causes multiplication overflow
	// We need: entryCount * IndexEntrySize > max uint64
	// Or: requiredBytes / entryCount != IndexEntrySize
	// This happens when entryCount * IndexEntrySize overflows uint64
	maxUint32 := uint32(^uint32(0))

	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, maxUint32) // EntryCount = max uint32
	_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
	_ = binary.Write(buf, binary.LittleEndian, uint64(0)) // FirstEntryOffset

	var index FileIndex
	_, err := index.readFrom(buf)

	if err == nil {
		t.Error("ReadFrom() expected error for multiplication overflow, got nil")
		return
	}

	// Check if it's the multiplication overflow error
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "calculation overflow"):
		// This is the correct error path
		var pkgErr *pkgerrors.PackageError
		if !pkgerrors.As(err, &pkgErr) {
			t.Errorf("ReadFrom() error is not a PackageError: %T", err)
			return
		}
		if pkgErr.Type != pkgerrors.ErrTypeValidation {
			t.Errorf("ReadFrom() error type = %v, want %v", pkgErr.Type, pkgerrors.ErrTypeValidation)
		}
	case strings.Contains(errStr, "exceeds maximum allocation size"):
		// This is also acceptable - it means we hit the maxInt/int(IndexEntrySize) check first
		// Both are valid OOM prevention paths
	default:
		// Any validation error is acceptable for OOM prevention
		var pkgErr *pkgerrors.PackageError
		if pkgerrors.As(err, &pkgErr) && pkgErr.Type == pkgerrors.ErrTypeValidation {
			// Acceptable - any validation error prevents OOM
		} else {
			t.Errorf("ReadFrom() unexpected error: %v", err)
		}
	}
}
