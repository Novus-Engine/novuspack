package metadata

import (
	"bytes"
	"testing"

	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

// TestOptionalDataEntry_WriteTo tests the WriteTo method.
func TestOptionalDataEntry_WriteTo(t *testing.T) {
	runWriteToEntryTable(t, []writeToCase{
		{"valid optional data entry", &OptionalDataEntry{
			DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10),
		}, false, 13},
		{"empty data", &OptionalDataEntry{
			DataType: OptionalDataTagsData, DataLength: 0, Data: []byte{},
		}, false, 3},
		{"data length mismatch", &OptionalDataEntry{
			DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 5),
		}, true, 0},
		{"incomplete write", &OptionalDataEntry{
			DataType: OptionalDataTagsData, DataLength: 10, Data: make([]byte, 10),
		}, false, 13},
	})
}

// TestOptionalDataEntry_WriteTo_ErrorPaths tests error paths in WriteTo method.
func TestOptionalDataEntry_WriteTo_ErrorPaths(t *testing.T) {
	entry := &OptionalDataEntry{
		DataType:   OptionalDataTagsData,
		DataLength: 10,
		Data:       make([]byte, 10),
	}
	runWriteToErrorPathsTable(t, entry, []writeToErrorCase{
		{"write error on DataType", testhelpers.NewErrorWriter(), true},
		{"write error on DataLength", testhelpers.NewFailingWriter(1), true},
		{"write error on Data", testhelpers.NewFailingWriter(3), true},
		{"incomplete write on Data", testhelpers.NewFailingWriter(5), true},
	})
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
	if _, err := entry.writeTo(&buf); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	var readEntry OptionalDataEntry
	n, err := readEntry.readFrom(&buf)

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

	if !bytes.Equal(readEntry.Data, entry.Data) {
		t.Errorf("ReadFrom() Data = %q, want %q", string(readEntry.Data), string(entry.Data))
	}
}

// TestOptionalDataEntry_ReadFrom_IncompleteData tests error handling for incomplete data.
func TestOptionalDataEntry_ReadFrom_IncompleteData(t *testing.T) {
	tests := []readFromIncompleteCase{
		{"no data", []byte{}},
		{"incomplete DataLength", []byte{0x00}},
		{"incomplete Data", []byte{0x00, 0x0A, 0x00}},
		{"partial Data read", []byte{0x00, 0x0A, 0x00, 0x74, 0x65, 0x73, 0x74}},
	}
	runReadFromIncompleteTable(t, tests, func() readFromEntry { return &OptionalDataEntry{} })
}
