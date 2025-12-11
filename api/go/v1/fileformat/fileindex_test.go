package fileformat

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"
)

// Helper types are defined in test_helpers.go

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
			err := tt.index.Validate()
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

			if index.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", index.Size(), tt.wantSize)
			}
		})
	}
}

// TestNewFileIndex verifies NewFileIndex initializes correctly
func TestNewFileIndex(t *testing.T) {
	index := NewFileIndex()

	if index == nil {
		t.Fatal("NewFileIndex() returned nil")
	}

	// Verify all fields are zero or empty
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
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
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
			FileIndex{
				EntryCount: 5,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024},
					{FileID: 5, Offset: 2048},
				},
			},
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
			n, err := index.ReadFrom(buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.index.Size()
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
				for i, entry := range index.Entries {
					if entry.FileID != tt.index.Entries[i].FileID {
						t.Errorf("Entry[%d].FileID = %d, want %d", i, entry.FileID, tt.index.Entries[i].FileID)
					}
					if entry.Offset != tt.index.Entries[i].Offset {
						t.Errorf("Entry[%d].Offset = %d, want %d", i, entry.Offset, tt.index.Entries[i].Offset)
					}
				}

				// Verify validation passes
				if err := index.Validate(); err != nil {
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
			_ = binary.Write(buf, binary.LittleEndian, uint32(2)) // EntryCount = 2
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
			_ = binary.Write(buf, binary.LittleEndian, uint64(0)) // FirstEntryOffset
			// Only 16 bytes, but EntryCount says 2 entries needed (32 more bytes)
			return buf.Bytes()
		}()},
		{"Header with EntryCount>0 but incomplete first entry", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2)) // EntryCount = 2
			_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // Reserved
			_ = binary.Write(buf, binary.LittleEndian, uint64(0)) // FirstEntryOffset
			_ = binary.Write(buf, binary.LittleEndian, uint64(1)) // First entry FileID
			// Only 8 bytes of first entry (need 16)
			return buf.Bytes()
		}()},
		{"Header with EntryCount>0 but incomplete second entry", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint32(2))   // EntryCount = 2
			_ = binary.Write(buf, binary.LittleEndian, uint32(0))   // Reserved
			_ = binary.Write(buf, binary.LittleEndian, uint64(0))   // FirstEntryOffset
			_ = binary.Write(buf, binary.LittleEndian, uint64(1))   // First entry FileID
			_ = binary.Write(buf, binary.LittleEndian, uint64(112)) // First entry Offset
			_ = binary.Write(buf, binary.LittleEndian, uint64(2))   // Second entry FileID
			// Only 8 bytes of second entry (need 16)
			return buf.Bytes()
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var index FileIndex
			r := bytes.NewReader(tt.data)
			_, err := index.ReadFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestFileIndexWriteTo verifies WriteTo serialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 5 - File Index Section
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
			FileIndex{
				EntryCount: 5,
				Reserved:   0,
				Entries: []IndexEntry{
					{FileID: 1, Offset: 112},
					{FileID: 2, Offset: 256},
					{FileID: 3, Offset: 512},
					{FileID: 4, Offset: 1024},
					{FileID: 5, Offset: 2048},
				},
			},
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
			n, err := tt.index.WriteTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.index.Size()
				if n != int64(expectedSize) {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, expectedSize)
				}

				if buf.Len() != expectedSize {
					t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), expectedSize)
				}

				// Verify we can read it back
				var index FileIndex
				_, readErr := index.ReadFrom(&buf)
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
			&errorWriter{},
			true,
			"failed to write",
		},
		{
			"Failing writer during Reserved write",
			FileIndex{
				EntryCount: 1,
				Entries:    []IndexEntry{{FileID: 1, Offset: 112}},
			},
			&failingWriter{maxBytes: 3},
			true,
			"failed to write",
		},
		{
			"Failing writer during FirstEntryOffset write",
			FileIndex{
				EntryCount: 1,
				Entries:    []IndexEntry{{FileID: 1, Offset: 112}},
			},
			&failingWriter{maxBytes: 7},
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
			&failingWriter{maxBytes: 20},
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
			&failingWriter{maxBytes: 32}, // Allow header (16) + first entry (16) but fail on second
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
			&failingWriter{maxBytes: 48}, // Allow header (16) + 2 entries (32) but fail on third
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
			&failingWriter{maxBytes: 64}, // Allow header (16) + 3 entries (48) but fail on fourth
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
			_, err := tt.index.WriteTo(tt.writer)

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
			if _, err := tt.index.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var index FileIndex
			if _, err := index.ReadFrom(&buf); err != nil {
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
			for i, entry := range index.Entries {
				if entry.FileID != tt.index.Entries[i].FileID {
					t.Errorf("Entry[%d].FileID mismatch: %d != %d", i, entry.FileID, tt.index.Entries[i].FileID)
				}
				if entry.Offset != tt.index.Entries[i].Offset {
					t.Errorf("Entry[%d].Offset mismatch: %d != %d", i, entry.Offset, tt.index.Entries[i].Offset)
				}
			}

			// Validate
			if err := index.Validate(); err != nil {
				t.Errorf("Round-trip index validation failed: %v", err)
			}

			// Verify FirstEntryOffset matches
			if index.FirstEntryOffset != tt.index.FirstEntryOffset {
				t.Errorf("FirstEntryOffset mismatch: %d != %d", index.FirstEntryOffset, tt.index.FirstEntryOffset)
			}
		})
	}
}
