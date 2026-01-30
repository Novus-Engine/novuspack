package metadata

import (
	"bytes"
	"io"
	"testing"

	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestOptionalDataEntry_WriteTo tests the WriteTo method.
func TestOptionalDataEntry_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		entry   OptionalDataEntry
		wantErr bool
	}{
		{
			name: "valid optional data entry",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			wantErr: false,
		},
		{
			name: "empty data",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 0,
				Data:       []byte{},
			},
			wantErr: false,
		},
		{
			name: "data length mismatch",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 5), // Mismatch
			},
			wantErr: true,
		},
		{
			name: "incomplete write",
			entry: OptionalDataEntry{
				DataType:   OptionalDataTagsData,
				DataLength: 10,
				Data:       make([]byte, 10),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			n, err := tt.entry.WriteTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if n == 0 {
					t.Error("WriteTo() wrote 0 bytes")
				}

				// Verify minimum size: DataType (1) + DataLength (2) + Data
				minSize := int64(3 + len(tt.entry.Data))
				if n < minSize {
					t.Errorf("WriteTo() wrote %d bytes, want at least %d", n, minSize)
				}
			}
		})
	}
}

// TestOptionalDataEntry_WriteTo_ErrorPaths tests error paths in WriteTo method.
func TestOptionalDataEntry_WriteTo_ErrorPaths(t *testing.T) {
	entry := OptionalDataEntry{
		DataType:   OptionalDataTagsData,
		DataLength: 10,
		Data:       make([]byte, 10),
	}

	tests := []struct {
		name    string
		writer  io.Writer
		wantErr bool
	}{
		{
			name:    "write error on DataType",
			writer:  testhelpers.NewErrorWriter(),
			wantErr: true,
		},
		{
			name:    "write error on DataLength",
			writer:  testhelpers.NewFailingWriter(1), // Fails after writing DataType
			wantErr: true,
		},
		{
			name:    "write error on Data",
			writer:  testhelpers.NewFailingWriter(3), // Fails after writing DataType and DataLength
			wantErr: true,
		},
		{
			name:    "incomplete write on Data",
			writer:  testhelpers.NewFailingWriter(5), // Fails after writing DataType, DataLength, and partial Data
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entry.WriteTo(tt.writer)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Error("WriteTo() expected error but got nil")
				}
				// Note: errorWriter returns error immediately, so bytes written may be 0
				// failingWriter may write some bytes before failing
			}
		})
	}
}

// TestOptionalDataEntry_ReadFrom tests the ReadFrom method.
func TestOptionalDataEntry_ReadFrom(t *testing.T) {
	// Create a valid entry to write
	entry := OptionalDataEntry{
		DataType:   OptionalDataTagsData,
		DataLength: 10,
		Data:       []byte("test data!"),
	}

	var buf bytes.Buffer
	if _, err := entry.WriteTo(&buf); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	var readEntry OptionalDataEntry
	n, err := readEntry.ReadFrom(&buf)

	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	if n == 0 {
		t.Error("ReadFrom() read 0 bytes")
	}

	if readEntry.DataType != entry.DataType {
		t.Errorf("ReadFrom() DataType = %v, want %v", readEntry.DataType, entry.DataType)
	}

	if readEntry.DataLength != entry.DataLength {
		t.Errorf("ReadFrom() DataLength = %v, want %v", readEntry.DataLength, entry.DataLength)
	}

	if string(readEntry.Data) != string(entry.Data) {
		t.Errorf("ReadFrom() Data = %q, want %q", string(readEntry.Data), string(entry.Data))
	}
}

// TestOptionalDataEntry_ReadFrom_IncompleteData tests error handling for incomplete data.
func TestOptionalDataEntry_ReadFrom_IncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "no data",
			data: []byte{},
		},
		{
			name: "incomplete DataLength",
			data: []byte{0x00}, // Only 1 byte (need 2 for DataLength)
		},
		{
			name: "incomplete Data",
			data: []byte{0x00, 0x0A, 0x00}, // DataType + DataLength (10) + only 1 byte of data (need 10)
		},
		{
			name: "partial Data read",
			data: []byte{0x00, 0x0A, 0x00, 0x74, 0x65, 0x73, 0x74}, // DataType + DataLength (10) + only 4 bytes of data (need 10)
		},
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
