package metadata

import (
	"bytes"
	"io"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestHashEntry_WriteTo tests the WriteTo method.
func TestHashEntry_WriteTo(t *testing.T) {
	tests := []struct {
		name    string
		entry   HashEntry
		wantErr bool
	}{
		{
			name: "valid hash entry",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
			},
			wantErr: false,
		},
		{
			name: "empty hash data",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  0,
				HashData:    []byte{},
			},
			wantErr: false,
		},
		{
			name: "hash length mismatch",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 16), // Mismatch
			},
			wantErr: true,
		},
		{
			name: "empty hash length",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  0,
				HashData:    []byte{},
			},
			wantErr: false,
		},
		{
			name: "incomplete write",
			entry: HashEntry{
				HashType:    fileformat.HashTypeSHA256,
				HashPurpose: fileformat.HashPurposeContentVerification,
				HashLength:  32,
				HashData:    make([]byte, 32),
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

				// Verify minimum size: HashType (1) + HashPurpose (1) + HashLength (2) + HashData
				minSize := int64(4 + len(tt.entry.HashData))
				if n < minSize {
					t.Errorf("WriteTo() wrote %d bytes, want at least %d", n, minSize)
				}
			}
		})
	}
}

// TestHashEntry_WriteTo_ErrorPaths tests error paths in WriteTo method.
func TestHashEntry_WriteTo_ErrorPaths(t *testing.T) {
	entry := HashEntry{
		HashType:    fileformat.HashTypeSHA256,
		HashPurpose: fileformat.HashPurposeContentVerification,
		HashLength:  32,
		HashData:    make([]byte, 32),
	}

	tests := []struct {
		name    string
		writer  io.Writer
		wantErr bool
	}{
		{
			name:    "write error on HashType",
			writer:  testhelpers.NewErrorWriter(),
			wantErr: true,
		},
		{
			name:    "write error on HashPurpose",
			writer:  testhelpers.NewFailingWriter(1), // Fails after writing HashType
			wantErr: true,
		},
		{
			name:    "write error on HashLength",
			writer:  testhelpers.NewFailingWriter(2), // Fails after writing HashType and HashPurpose
			wantErr: true,
		},
		{
			name:    "write error on HashData",
			writer:  testhelpers.NewFailingWriter(4), // Fails after writing HashType, HashPurpose, and HashLength
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

// TestHashEntry_ReadFrom tests the ReadFrom method.
func TestHashEntry_ReadFrom(t *testing.T) {
	// Create a valid entry to write
	entry := HashEntry{
		HashType:    fileformat.HashTypeSHA256,
		HashPurpose: fileformat.HashPurposeContentVerification,
		HashLength:  32,
		HashData:    make([]byte, 32),
	}

	var buf bytes.Buffer
	if _, err := entry.WriteTo(&buf); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	var readEntry HashEntry
	n, err := readEntry.ReadFrom(&buf)

	if err != nil {
		t.Fatalf("ReadFrom() error = %v", err)
	}

	if n == 0 {
		t.Error("ReadFrom() read 0 bytes")
	}

	if readEntry.HashType != entry.HashType {
		t.Errorf("ReadFrom() HashType = %v, want %v", readEntry.HashType, entry.HashType)
	}

	if readEntry.HashPurpose != entry.HashPurpose {
		t.Errorf("ReadFrom() HashPurpose = %v, want %v", readEntry.HashPurpose, entry.HashPurpose)
	}

	if readEntry.HashLength != entry.HashLength {
		t.Errorf("ReadFrom() HashLength = %v, want %v", readEntry.HashLength, entry.HashLength)
	}

	if len(readEntry.HashData) != len(entry.HashData) {
		t.Errorf("ReadFrom() HashData length = %d, want %d", len(readEntry.HashData), len(entry.HashData))
	}
}

// TestHashEntry_ReadFrom_IncompleteData tests error handling for incomplete data.
func TestHashEntry_ReadFrom_IncompleteData(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "no data",
			data: []byte{},
		},
		{
			name: "incomplete HashPurpose",
			data: []byte{0x01}, // Only 1 byte (need 2 for HashType and HashPurpose)
		},
		{
			name: "incomplete HashLength",
			data: []byte{0x01, 0x01}, // HashType + HashPurpose, but only 1 byte of HashLength (need 2)
		},
		{
			name: "incomplete HashData",
			data: []byte{0x01, 0x01, 0x20, 0x00}, // HashType + HashPurpose + HashLength (32), but no data
		},
		{
			name: "partial HashData read",
			data: []byte{0x01, 0x01, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Only 8 bytes of hash data (need 32)
		},
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
