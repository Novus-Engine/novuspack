package metadata

import (
	"bytes"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestHashEntry_WriteTo tests the WriteTo method.
func TestHashEntry_WriteTo(t *testing.T) {
	runWriteToEntryTable(t, []writeToCase{
		{"valid hash entry", &HashEntry{
			HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification,
			HashLength: 32, HashData: make([]byte, 32),
		}, false, 36},
		{"empty hash data", &HashEntry{
			HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification,
			HashLength: 0, HashData: []byte{},
		}, false, 4},
		{"hash length mismatch", &HashEntry{
			HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification,
			HashLength: 32, HashData: make([]byte, 16),
		}, true, 0},
		{"empty hash length", &HashEntry{
			HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification,
			HashLength: 0, HashData: []byte{},
		}, false, 4},
		{"incomplete write", &HashEntry{
			HashType: fileformat.HashTypeSHA256, HashPurpose: fileformat.HashPurposeContentVerification,
			HashLength: 32, HashData: make([]byte, 32),
		}, false, 36},
	})
}

// TestHashEntry_WriteTo_ErrorPaths tests error paths in WriteTo method.
func TestHashEntry_WriteTo_ErrorPaths(t *testing.T) {
	entry := &HashEntry{
		HashType:    fileformat.HashTypeSHA256,
		HashPurpose: fileformat.HashPurposeContentVerification,
		HashLength:  32,
		HashData:    make([]byte, 32),
	}
	runWriteToErrorPathsTable(t, entry, []writeToErrorCase{
		{"write error on HashType", testhelpers.NewErrorWriter(), true},
		{"write error on HashPurpose", testhelpers.NewFailingWriter(1), true},
		{"write error on HashLength", testhelpers.NewFailingWriter(2), true},
		{"write error on HashData", testhelpers.NewFailingWriter(4), true},
	})
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
	if _, err := entry.writeTo(&buf); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	var readEntry HashEntry
	n, err := readEntry.readFrom(&buf)

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
	tests := []readFromIncompleteCase{
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

	runReadFromIncompleteTable(t, tests, func() readFromEntry { return &HashEntry{} })
}
