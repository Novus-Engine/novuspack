package testutil

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
)

// TestCreateTestPackageFile tests the CreateTestPackageFile helper function.
func TestCreateTestPackageFile(t *testing.T) {
	tmpDir := t.TempDir()
	pkgPath := filepath.Join(tmpDir, "test.nvpk")

	// Create test package file
	CreateTestPackageFile(t, pkgPath)

	// Verify file was created
	info, err := os.Stat(pkgPath)
	if err != nil {
		t.Fatalf("Failed to stat created file: %v", err)
	}

	if info.Size() == 0 {
		t.Error("Created file should not be empty")
	}

	// Verify file can be opened and read
	file, err := os.Open(pkgPath)
	if err != nil {
		t.Fatalf("Failed to open created file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// Read and validate header
	header := fileformat.NewPackageHeader()
	if err := binary.Read(file, binary.LittleEndian, header); err != nil {
		t.Fatalf("Failed to read header from created file: %v", err)
	}

	// Verify header has correct magic number
	if header.Magic != fileformat.NVPKMagic {
		t.Errorf("Header magic = 0x%08X, want 0x%08X", header.Magic, fileformat.NVPKMagic)
	}

	// Read and validate file index
	index := fileformat.NewFileIndex()
	if err := binary.Read(file, binary.LittleEndian, &index.EntryCount); err != nil {
		t.Fatalf("Failed to read index entry count from created file: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &index.Reserved); err != nil {
		t.Fatalf("Failed to read index reserved from created file: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &index.FirstEntryOffset); err != nil {
		t.Fatalf("Failed to read index first entry offset from created file: %v", err)
	}
	if index.EntryCount > 0 {
		index.Entries = make([]fileformat.IndexEntry, 0, index.EntryCount)
		for i := uint32(0); i < index.EntryCount; i++ {
			var entry fileformat.IndexEntry
			if err := binary.Read(file, binary.LittleEndian, &entry); err != nil {
				t.Fatalf("Failed to read index entry %d from created file: %v", i, err)
			}
			index.Entries = append(index.Entries, entry)
		}
	}

	// Verify index has zero entries (empty package)
	if index.EntryCount != 0 {
		t.Errorf("Index EntryCount = %d, want 0", index.EntryCount)
	}
}
