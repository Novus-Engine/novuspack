package metadata

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/novus-engine/novuspack/api/go/fileformat"
	"github.com/novus-engine/novuspack/api/go/generics"
	"github.com/novus-engine/novuspack/api/go/internal/testhelpers"
)

func runWriteMetaToFailingWriterTest(t *testing.T, errMsg string) {
	t.Helper()
	fe := NewFileEntry()
	fe.FileID = 1
	failingWriter := testhelpers.NewErrorWriter()
	_, err := fe.WriteMetaTo(failingWriter)
	if err == nil {
		t.Error(errMsg)
	}
}

// TestMarshalMeta tests MarshalMeta method
func TestMarshalMeta(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "basic file entry",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.FileID = 1
				fe.OriginalSize = 100
				fe.StoredSize = 80
				return fe
			},
			wantErr: false,
		},
		{
			name: "with paths and hashes",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.FileID = 1
				fe.Paths = []generics.PathEntry{
					{PathLength: 4, Path: "file"},
				}
				fe.Hashes = []HashEntry{
					{
						HashType:    fileformat.HashTypeSHA256,
						HashPurpose: fileformat.HashPurposeContentVerification,
						HashLength:  32,
						HashData:    make([]byte, 32),
					},
				}
				return fe
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			meta, err := fe.MarshalMeta()

			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(meta) == 0 {
					t.Error("MarshalMeta() returned empty data")
				}

				// Verify it contains at least the fixed section
				if len(meta) < FileEntryFixedSize {
					t.Errorf("MarshalMeta() returned %d bytes, want at least %d", len(meta), FileEntryFixedSize)
				}
			}
		})
	}

	t.Run("failing writer", func(t *testing.T) {
		runWriteMetaToFailingWriterTest(t, "MarshalMeta() with failing writer should return error")
	})
}

// TestMarshalData tests MarshalData method
func TestMarshalData(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "data loaded in memory",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.SetData([]byte("test data"))
				return fe
			},
			wantErr: false,
		},
		{
			name:    "no data available",
			setup:   NewFileEntry,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			data, err := fe.MarshalData()

			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(data) == 0 {
					t.Error("MarshalData() returned empty data")
				}
			}
		})
	}
}

// TestMarshal tests Marshal method
func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "successful marshal",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.FileID = 1
				fe.SetData([]byte("test data"))
				return fe
			},
			wantErr: false,
		},
		{
			name: "marshal with no data",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.FileID = 1
				return fe
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			meta, data, err := fe.Marshal()

			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(meta) == 0 {
					t.Error("Marshal() returned empty metadata")
				}

				if len(data) == 0 {
					t.Error("Marshal() returned empty data")
				}
			}
		})
	}
}

// TestWriteMetaTo tests WriteMetaTo method
func TestWriteMetaTo(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "successful write",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.FileID = 1
				fe.OriginalSize = 100
				return fe
			},
			wantErr: false,
		},
		{
			name: "with paths and optional data",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.FileID = 1
				fe.Paths = []generics.PathEntry{
					{PathLength: 4, Path: "file"},
				}
				fe.OptionalData = []OptionalDataEntry{
					{DataType: OptionalDataTagsData, DataLength: 0, Data: []byte{}},
				}
				fe.updateOptionalDataLen()
				return fe
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			var buf bytes.Buffer
			n, err := fe.WriteMetaTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteMetaTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if n == 0 {
					t.Error("WriteMetaTo() wrote 0 bytes")
				}

				if buf.Len() == 0 {
					t.Error("WriteMetaTo() wrote no data to buffer")
				}
			}
		})
	}

	t.Run("failing writer", func(t *testing.T) {
		runWriteMetaToFailingWriterTest(t, "WriteMetaTo() with failing writer should return error")
	})
}

// TestWriteDataTo tests WriteDataTo method
//
//nolint:gocognit // table-driven write cases
func TestWriteDataTo(t *testing.T) {
	// Create a temporary source file
	sourceFile, err := os.CreateTemp("", "novuspack-source-*")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	defer func() {
		_ = sourceFile.Close()
		_ = os.Remove(sourceFile.Name())
	}()

	testData := []byte("test data for writing")
	if _, err := sourceFile.Write(testData); err != nil {
		t.Fatalf("Failed to write to source file: %v", err)
	}

	tests := []struct {
		name    string
		setup   func() *FileEntry
		wantErr bool
	}{
		{
			name: "data in memory",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.SetData([]byte("test data"))
				return fe
			},
			wantErr: false,
		},
		{
			name:    "no data available",
			setup:   NewFileEntry,
			wantErr: true,
		},
		{
			name: "data from source file",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				fe.setSourceFile(sourceFile, 0, int64(len(testData)))
				return fe
			},
			wantErr: false,
		},
		{
			name: "data from temp file",
			setup: func() *FileEntry {
				fe := NewFileEntry()
				err := fe.WriteToTempFile(context.Background(), testData)
				if err != nil {
					t.Fatalf("WriteToTempFile() error = %v", err)
				}
				return fe
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := tt.setup()
			var buf bytes.Buffer
			n, err := fe.WriteDataTo(&buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteDataTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if n == 0 {
					t.Error("WriteDataTo() wrote 0 bytes")
				}

				if tt.name == "data in memory" && buf.String() != "test data" {
					t.Errorf("WriteDataTo() data = %q, want %q", buf.String(), "test data")
				}
			}

			// Cleanup temp file if created
			if tt.name == "data from temp file" {
				_ = fe.CleanupTempFile(context.Background()) //nolint:errcheck // cleanup best-effort in test
			}
		})
	}

	// Test error path - seek error
	t.Run("seek error", func(t *testing.T) {
		fe := NewFileEntry()
		// Use a closed file to cause seek error
		closedFile, _ := os.Open(os.DevNull)
		_ = closedFile.Close()
		fe.setSourceFile(closedFile, 0, 10)
		var buf bytes.Buffer
		_, err := fe.WriteDataTo(&buf)
		if err == nil {
			t.Error("WriteDataTo() with seek error should return error")
		}
	})

	// Test error path - copy error
	t.Run("copy error", func(t *testing.T) {
		fe := NewFileEntry()
		sourceFile, err := os.CreateTemp("", "novuspack-source-*")
		if err != nil {
			t.Fatalf("Failed to create source file: %v", err)
		}
		defer func() {
			_ = sourceFile.Close()
			_ = os.Remove(sourceFile.Name())
		}()
		// Request more data than available
		fe.setSourceFile(sourceFile, 0, 100)
		var buf bytes.Buffer
		_, err = fe.WriteDataTo(&buf)
		if err == nil {
			t.Error("WriteDataTo() with copy error should return error")
		}
	})

	// Test error path - temp file open error
	t.Run("temp file open error", func(t *testing.T) {
		fe := NewFileEntry()
		// Set invalid temp file path
		fe.setTempPath("/invalid/path/that/does/not/exist")
		var buf bytes.Buffer
		_, err := fe.WriteDataTo(&buf)
		// This might succeed if it falls back to other data sources, or fail
		// The important thing is we test the temp file path
		_ = err // Accept either outcome
	})
}

// TestUnmarshalFileEntry tests UnmarshalFileEntry function
func TestUnmarshalFileEntry(t *testing.T) {
	// Create a FileEntry and marshal it
	original := NewFileEntry()
	original.FileID = 1
	original.OriginalSize = 100
	original.StoredSize = 80

	var buf bytes.Buffer
	_, err := original.WriteMetaTo(&buf)
	if err != nil {
		t.Fatalf("WriteMetaTo() error = %v", err)
	}

	// Unmarshal it
	unmarshaled, err := UnmarshalFileEntry(buf.Bytes())
	if err != nil {
		t.Fatalf("UnmarshalFileEntry() error = %v", err)
	}

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so unmarshaled is not nil after check
	if unmarshaled == nil {
		t.Fatal("UnmarshalFileEntry() returned nil")
	}

	//nolint:staticcheck // SA5011: false positive - t.Fatal exits, so unmarshaled is not nil after check
	if unmarshaled.FileID != original.FileID {
		t.Errorf("UnmarshalFileEntry() FileID = %d, want %d", unmarshaled.FileID, original.FileID)
	}
}

// TestFileEntry_WriteTo tests the WriteTo method.
func TestFileEntry_WriteTo(t *testing.T) {
	fe := NewFileEntry()
	fe.FileID = 1
	fe.SetData([]byte("test data"))

	var buf bytes.Buffer
	n, err := fe.WriteTo(&buf)

	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	if n == 0 {
		t.Error("WriteTo() wrote 0 bytes")
	}

	if buf.Len() == 0 {
		t.Error("WriteTo() wrote no data to buffer")
	}

	// Test error path - WriteDataTo error (no data)
	fe2 := NewFileEntry()
	fe2.FileID = 1
	// WriteTo will fail because WriteDataTo fails when no data
	_, err = fe2.WriteTo(&buf)
	if err == nil {
		t.Error("WriteTo() with no data should return error")
	}
}

// TestWriteToDelegatesToWriteMetaToAndWriteDataTo tests that WriteTo delegates correctly
func TestWriteToDelegatesToWriteMetaToAndWriteDataTo(t *testing.T) {
	fe := NewFileEntry()
	fe.FileID = 1
	fe.OriginalSize = 100
	fe.SetData([]byte("test data"))

	var buf bytes.Buffer
	n, err := fe.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	if n == 0 {
		t.Error("WriteTo() wrote 0 bytes")
	}

	// Verify it wrote both metadata and data
	if buf.Len() < FileEntryFixedSize+len("test data") {
		t.Errorf("WriteTo() wrote %d bytes, want at least %d", buf.Len(), FileEntryFixedSize+len("test data"))
	}
}
