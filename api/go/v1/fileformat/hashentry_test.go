package fileformat

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"
)

// TestHashEntrySize verifies HashEntry size calculation
func TestHashEntrySize(t *testing.T) {
	tests := []struct {
		name       string
		hashLength uint16
		wantSize   int
	}{
		{"SHA256", 32, 36},
		{"SHA512", 64, 68},
		{"XXH3", 8, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := HashEntry{
				HashType:    HashTypeSHA256,
				HashPurpose: HashPurposeContentVerification,
				HashLength:  tt.hashLength,
				HashData:    make([]byte, tt.hashLength),
			}

			if entry.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", entry.Size(), tt.wantSize)
			}
		})
	}
}

// TestHashEntryValidation verifies validation logic
func TestHashEntryValidation(t *testing.T) {
	tests := []struct {
		name    string
		entry   HashEntry
		wantErr bool
	}{
		{
			"Valid SHA256",
			HashEntry{
				HashType:    HashTypeSHA256,
				HashPurpose: HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
			false,
		},
		{
			"Nil hash data",
			HashEntry{
				HashType:   HashTypeSHA256,
				HashLength: 32,
				HashData:   nil,
			},
			true,
		},
		{
			"Length mismatch",
			HashEntry{
				HashType:   HashTypeSHA256,
				HashLength: 32,
				HashData:   make([]byte, 16),
			},
			true,
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

// TestHashEntryReadFrom verifies ReadFrom deserialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3 - Hash Data
func TestHashEntryReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		entry   HashEntry
		wantErr bool
	}{
		{
			"Valid SHA256 entry",
			HashEntry{
				HashType:    HashTypeSHA256,
				HashPurpose: HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
			false,
		},
		{
			"Valid SHA512 entry",
			HashEntry{
				HashType:    HashTypeSHA512,
				HashPurpose: HashPurposeDeduplication,
				HashLength:  64,
				HashData:    make([]byte, 64),
			},
			false,
		},
		{
			"Valid XXH3 entry",
			HashEntry{
				HashType:    HashTypeXXH3,
				HashPurpose: HashPurposeFastLookup,
				HashLength:  8,
				HashData:    make([]byte, 8),
			},
			false,
		},
		{
			"Empty hash entry (HashLength = 0)",
			HashEntry{
				HashType:    HashTypeSHA256,
				HashPurpose: HashPurposeContentVerification,
				HashLength:  0,
				HashData:    nil,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure HashLength matches actual data length
			tt.entry.HashLength = uint16(len(tt.entry.HashData))

			// Serialize original entry using WriteTo (once implemented)
			// For now, serialize manually
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, tt.entry.HashType)
			_ = binary.Write(buf, binary.LittleEndian, tt.entry.HashPurpose)
			_ = binary.Write(buf, binary.LittleEndian, tt.entry.HashLength)
			buf.Write(tt.entry.HashData)

			// Deserialize using ReadFrom
			var entry HashEntry
			n, err := entry.ReadFrom(buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.entry.Size()
				if n != int64(expectedSize) {
					t.Errorf("ReadFrom() read %d bytes, want %d", n, expectedSize)
				}

				// Verify all fields match
				if entry.HashType != tt.entry.HashType {
					t.Errorf("HashType = %d, want %d", entry.HashType, tt.entry.HashType)
				}
				if entry.HashPurpose != tt.entry.HashPurpose {
					t.Errorf("HashPurpose = %d, want %d", entry.HashPurpose, tt.entry.HashPurpose)
				}
				if entry.HashLength != tt.entry.HashLength {
					t.Errorf("HashLength = %d, want %d", entry.HashLength, tt.entry.HashLength)
				}
				if len(entry.HashData) != len(tt.entry.HashData) {
					t.Errorf("HashData length = %d, want %d", len(entry.HashData), len(tt.entry.HashData))
				}

				// Verify validation passes (skip for empty entries as they're invalid)
				if tt.entry.HashLength > 0 {
					if err := entry.Validate(); err != nil {
						t.Errorf("ReadFrom() entry validation failed: %v", err)
					}
				}
			}
		})
	}
}

// TestHashEntryReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestHashEntryReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Only HashType", []byte{0x00}},
		{"Missing HashLength", []byte{0x00, 0x00}},
		{"Incomplete HashLength read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0)) // HashType
			_ = binary.Write(buf, binary.LittleEndian, uint8(0)) // HashPurpose
			return buf.Bytes()[:3]                               // Only 1 byte of HashLength (need 2)
		}()},
		{"Incomplete hash data", []byte{0x00, 0x00, 0x20, 0x00, 0x01, 0x02, 0x03}}, // Only 3 bytes of 32
		{"HashLength > 0 but no data", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // HashType
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // HashPurpose
			_ = binary.Write(buf, binary.LittleEndian, uint16(32)) // HashLength = 32
			// No hash data
			return buf.Bytes()
		}()},
		{"HashLength > 0 but partial data", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // HashType
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // HashPurpose
			_ = binary.Write(buf, binary.LittleEndian, uint16(32)) // HashLength = 32
			partialData := make([]byte, 15)                        // Only 15 bytes of 32
			buf.Write(partialData)
			return buf.Bytes()
		}()},
		{"HashLength > 0 but partial data for large hash", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(1))   // HashType (SHA512)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // HashPurpose
			_ = binary.Write(buf, binary.LittleEndian, uint16(64)) // HashLength = 64
			partialData := make([]byte, 30)                        // Only 30 bytes of 64
			buf.Write(partialData)
			return buf.Bytes()
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry HashEntry
			r := bytes.NewReader(tt.data)
			_, err := entry.ReadFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestHashEntryWriteTo verifies WriteTo serialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.3 - Hash Data
func TestHashEntryWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		entry   HashEntry
		wantErr bool
	}{
		{
			"Valid SHA256 entry",
			HashEntry{
				HashType:    HashTypeSHA256,
				HashPurpose: HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
			false,
		},
		{
			"Valid SHA512 entry",
			HashEntry{
				HashType:    HashTypeSHA512,
				HashPurpose: HashPurposeDeduplication,
				HashLength:  64,
				HashData:    make([]byte, 64),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure HashLength matches actual data length
			tt.entry.HashLength = uint16(len(tt.entry.HashData))

			var buf bytes.Buffer
			n, err := tt.entry.WriteTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				expectedSize := tt.entry.Size()
				if n != int64(expectedSize) {
					t.Errorf("WriteTo() wrote %d bytes, want %d", n, expectedSize)
				}

				if buf.Len() != expectedSize {
					t.Errorf("WriteTo() buffer size = %d bytes, want %d", buf.Len(), expectedSize)
				}

				// Verify we can read it back
				var entry HashEntry
				_, readErr := entry.ReadFrom(&buf)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if entry.HashType != tt.entry.HashType {
					t.Errorf("HashType mismatch: %d != %d", entry.HashType, tt.entry.HashType)
				}
				if entry.HashLength != tt.entry.HashLength {
					t.Errorf("HashLength mismatch: %d != %d", entry.HashLength, tt.entry.HashLength)
				}
			}
		})
	}
}

// TestHashEntryRoundTrip verifies round-trip serialization
func TestHashEntryRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		entry HashEntry
	}{
		{
			"SHA256 entry",
			HashEntry{
				HashType:    HashTypeSHA256,
				HashPurpose: HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
		},
		{
			"SHA512 entry",
			HashEntry{
				HashType:    HashTypeSHA512,
				HashPurpose: HashPurposeDeduplication,
				HashLength:  64,
				HashData:    make([]byte, 64),
			},
		},
		{
			"BLAKE3 entry",
			HashEntry{
				HashType:    HashTypeBLAKE3,
				HashPurpose: HashPurposeIntegrity,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure HashLength matches actual data length
			tt.entry.HashLength = uint16(len(tt.entry.HashData))

			// Write
			var buf bytes.Buffer
			if _, err := tt.entry.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var entry HashEntry
			if _, err := entry.ReadFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			// Compare all fields
			if entry.HashType != tt.entry.HashType {
				t.Errorf("HashType mismatch: %d != %d", entry.HashType, tt.entry.HashType)
			}
			if entry.HashPurpose != tt.entry.HashPurpose {
				t.Errorf("HashPurpose mismatch: %d != %d", entry.HashPurpose, tt.entry.HashPurpose)
			}
			if entry.HashLength != tt.entry.HashLength {
				t.Errorf("HashLength mismatch: %d != %d", entry.HashLength, tt.entry.HashLength)
			}
			if len(entry.HashData) != len(tt.entry.HashData) {
				t.Errorf("HashData length mismatch: %d != %d", len(entry.HashData), len(tt.entry.HashData))
			}

			// Validate
			if err := entry.Validate(); err != nil {
				t.Errorf("Round-trip entry validation failed: %v", err)
			}
		})
	}
}

// TestHashEntryWriteToErrorPaths verifies WriteTo error handling
func TestHashEntryWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		entry     HashEntry
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during HashType write",
			HashEntry{
				HashType:   HashTypeSHA256,
				HashLength: 32,
				HashData:   make([]byte, 32),
			},
			&errorWriter{},
			true,
			"failed to write",
		},
		{
			"Failing writer during hash data write",
			HashEntry{
				HashType:   HashTypeSHA256,
				HashLength: 32,
				HashData:   make([]byte, 32),
			},
			&failingWriter{maxBytes: 3},
			true,
			"failed to write",
		},
		{
			"Incomplete hash data write",
			HashEntry{
				HashType:   HashTypeSHA256,
				HashLength: 32,
				HashData:   make([]byte, 32),
			},
			&incompleteWriter{maxWrite: 20},
			true,
			"incomplete hash data write",
		},
		{
			"Hash length mismatch",
			HashEntry{
				HashType:   HashTypeSHA256,
				HashLength: 64,
				HashData:   make([]byte, 32), // 32 bytes, but HashLength says 64
			},
			&bytes.Buffer{},
			true,
			"hash length mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only set HashLength if not testing mismatch
			if !strings.Contains(tt.name, "mismatch") {
				tt.entry.HashLength = uint16(len(tt.entry.HashData))
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
