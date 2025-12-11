package fileformat

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"
)

// TestOptionalDataEntrySize verifies size calculation
func TestOptionalDataEntrySize(t *testing.T) {
	tests := []struct {
		name       string
		dataLength uint16
		wantSize   int
	}{
		{"Small data", 4, 7},
		{"Medium data", 64, 67},
		{"Large data", 1024, 1027},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: tt.dataLength,
				Data:       make([]byte, tt.dataLength),
			}

			if entry.Size() != tt.wantSize {
				t.Errorf("Size() = %d, want %d", entry.Size(), tt.wantSize)
			}
		})
	}
}

// TestOptionalDataEntryValidation verifies validation logic
func TestOptionalDataEntryValidation(t *testing.T) {
	tests := []struct {
		name    string
		entry   OptionalDataEntry
		wantErr bool
	}{
		{
			"Valid entry",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 4,
				Data:       []byte{1, 2, 3, 4},
			},
			false,
		},
		{
			"Nil data",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 4,
				Data:       nil,
			},
			true,
		},
		{
			"Length mismatch",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 4,
				Data:       []byte{1, 2},
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

// TestOptionalDataEntryReadFrom verifies ReadFrom deserialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4 - Optional Data
func TestOptionalDataEntryReadFrom(t *testing.T) {
	tests := []struct {
		name    string
		entry   OptionalDataEntry
		wantErr bool
	}{
		{
			"Valid entry with small data",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			false,
		},
		{
			"Valid entry with large data",
			OptionalDataEntry{
				DataType:   OptionalDataExtendedAttributes,
				DataLength: 100,
				Data:       make([]byte, 100),
			},
			false,
		},
		{
			"Valid entry with Windows attributes",
			OptionalDataEntry{
				DataType:   OptionalDataWindowsAttributes,
				DataLength: 4,
				Data:       make([]byte, 4),
			},
			false,
		},
		{
			"Empty data entry (DataLength = 0)",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 0,
				Data:       nil,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure DataLength matches actual data length
			tt.entry.DataLength = uint16(len(tt.entry.Data))

			// Serialize original entry manually
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, tt.entry.DataType)
			_ = binary.Write(buf, binary.LittleEndian, tt.entry.DataLength)
			buf.Write(tt.entry.Data)

			// Deserialize using ReadFrom
			var entry OptionalDataEntry
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
				if entry.DataType != tt.entry.DataType {
					t.Errorf("DataType = %d, want %d", entry.DataType, tt.entry.DataType)
				}
				if entry.DataLength != tt.entry.DataLength {
					t.Errorf("DataLength = %d, want %d", entry.DataLength, tt.entry.DataLength)
				}
				if len(entry.Data) != len(tt.entry.Data) {
					t.Errorf("Data length = %d, want %d", len(entry.Data), len(tt.entry.Data))
				}

				// Verify validation passes (skip for empty entries as they're invalid)
				if tt.entry.DataLength > 0 {
					if err := entry.Validate(); err != nil {
						t.Errorf("ReadFrom() entry validation failed: %v", err)
					}
				}
			}
		})
	}
}

// TestOptionalDataEntryReadFromIncompleteData verifies ReadFrom handles incomplete data
func TestOptionalDataEntryReadFromIncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"No data", []byte{}},
		{"Only DataType", []byte{0x00}},
		{"Incomplete DataLength read", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0)) // DataType
			return buf.Bytes()[:2]                               // Only 1 byte of DataLength (need 2)
		}()},
		{"Missing DataLength", []byte{0x00, 0x0A}},
		{"DataLength > 0 but no data", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // DataType
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // DataLength = 10
			// No data
			return buf.Bytes()
		}()},
		{"Incomplete data", []byte{0x00, 0x0A, 0x00, 0x01, 0x02}}, // Only 3 bytes of 10
		{"DataLength > 0 but partial data", func() []byte {
			buf := new(bytes.Buffer)
			_ = binary.Write(buf, binary.LittleEndian, uint8(0))   // DataType
			_ = binary.Write(buf, binary.LittleEndian, uint16(10)) // DataLength = 10
			partialData := make([]byte, 5)                         // Only 5 bytes of 10
			buf.Write(partialData)
			return buf.Bytes()
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var entry OptionalDataEntry
			r := bytes.NewReader(tt.data)
			_, err := entry.ReadFrom(r)

			if err == nil {
				t.Errorf("ReadFrom() expected error for incomplete data, got nil")
			}
		})
	}
}

// TestOptionalDataEntryWriteTo verifies WriteTo serialization
// Specification: ../../docs/tech_specs/package_file_format.md Section 4.1.4.4 - Optional Data
func TestOptionalDataEntryWriteTo(t *testing.T) {
	tests := []struct {
		name    string
		entry   OptionalDataEntry
		wantErr bool
	}{
		{
			"Valid entry with small data",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			false,
		},
		{
			"Valid entry with large data",
			OptionalDataEntry{
				DataType:   OptionalDataExtendedAttributes,
				DataLength: 100,
				Data:       make([]byte, 100),
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure DataLength matches actual data length
			tt.entry.DataLength = uint16(len(tt.entry.Data))

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
				var entry OptionalDataEntry
				_, readErr := entry.ReadFrom(&buf)
				if readErr != nil {
					t.Errorf("Failed to read back written data: %v", readErr)
				}

				if entry.DataType != tt.entry.DataType {
					t.Errorf("DataType mismatch: %d != %d", entry.DataType, tt.entry.DataType)
				}
				if entry.DataLength != tt.entry.DataLength {
					t.Errorf("DataLength mismatch: %d != %d", entry.DataLength, tt.entry.DataLength)
				}
			}
		})
	}
}

// TestOptionalDataEntryRoundTrip verifies round-trip serialization
func TestOptionalDataEntryRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		entry OptionalDataEntry
	}{
		{
			"Small data entry",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
		},
		{
			"Large data entry",
			OptionalDataEntry{
				DataType:   OptionalDataExtendedAttributes,
				DataLength: 200,
				Data:       make([]byte, 200),
			},
		},
		{
			"Windows attributes entry",
			OptionalDataEntry{
				DataType:   OptionalDataWindowsAttributes,
				DataLength: 4,
				Data:       make([]byte, 4),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure DataLength matches actual data length
			tt.entry.DataLength = uint16(len(tt.entry.Data))

			// Write
			var buf bytes.Buffer
			if _, err := tt.entry.WriteTo(&buf); err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			// Read
			var entry OptionalDataEntry
			if _, err := entry.ReadFrom(&buf); err != nil {
				t.Fatalf("ReadFrom() error = %v", err)
			}

			// Compare all fields
			if entry.DataType != tt.entry.DataType {
				t.Errorf("DataType mismatch: %d != %d", entry.DataType, tt.entry.DataType)
			}
			if entry.DataLength != tt.entry.DataLength {
				t.Errorf("DataLength mismatch: %d != %d", entry.DataLength, tt.entry.DataLength)
			}
			if len(entry.Data) != len(tt.entry.Data) {
				t.Errorf("Data length mismatch: %d != %d", len(entry.Data), len(tt.entry.Data))
			}

			// Validate
			if err := entry.Validate(); err != nil {
				t.Errorf("Round-trip entry validation failed: %v", err)
			}
		})
	}
}

// TestOptionalDataEntryWriteToErrorPaths verifies WriteTo error handling
func TestOptionalDataEntryWriteToErrorPaths(t *testing.T) {
	tests := []struct {
		name      string
		entry     OptionalDataEntry
		writer    io.Writer
		wantErr   bool
		errSubstr string
	}{
		{
			"Error writer during DataType write",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			&errorWriter{},
			true,
			"failed to write",
		},
		{
			"Failing writer during data write",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			&failingWriter{maxBytes: 2},
			true,
			"failed to write",
		},
		{
			"Incomplete data write",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			&incompleteWriter{maxWrite: 5},
			true,
			"incomplete data write",
		},
		{
			"Data length mismatch",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 20,
				Data:       make([]byte, 10), // 10 bytes, but DataLength says 20
			},
			&bytes.Buffer{},
			true,
			"data length mismatch",
		},
		{
			"Empty data entry",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 0,
				Data:       nil,
			},
			&bytes.Buffer{},
			false,
			"",
		},
		{
			"Failing writer during large data write",
			OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 100,
				Data:       make([]byte, 100),
			},
			&failingWriter{maxBytes: 10},
			true,
			"failed to write",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only set DataLength if not testing mismatch
			if !strings.Contains(tt.name, "mismatch") {
				tt.entry.DataLength = uint16(len(tt.entry.Data))
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
